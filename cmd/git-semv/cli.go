package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	flags "github.com/jessevdk/go-flags"
	semv "github.com/linyows/git-semv/git"
)

const (
	// ExitOK for exit code
	ExitOK int = 0

	// ExitErr for exit code
	ExitErr int = 1
)

// CLI struct
type CLI struct {
	outStream, errStream io.Writer
	Command              string
	Args                 []string
	Major                bool   `long:"major" short:"M" description:"Major version when you make incompatible API changes"`
	Minor                bool   `long:"minor" short:"m" description:"Minor version when you add functionality in a backwards-compatible manner: default"`
	Patch                bool   `long:"patch" short:"p" description:"Patch version when you make backwards-compatible bug fixes"`
	Pre                  bool   `long:"pre" short:"P" description:"Pre-Release version indicates"`
	Prefix               string `long:"prefix" short:"x" description:"Prefix for version and tag(default: v)"`
	Help                 bool   `long:"help" short:"h" description:"Show this help message and exit"`
	Version              bool   `long:"version" short:"v" description:"Prints the version number"`
}

func (c *CLI) buildHelp(names []string) []string {
	var help []string
	t := reflect.TypeOf(CLI{})

	for _, name := range names {
		f, ok := t.FieldByName(name)
		if !ok {
			continue
		}

		tag := f.Tag
		if tag == "" {
			continue
		}

		var o, a string
		if a = tag.Get("arg"); a != "" {
			a = fmt.Sprintf("=%s", a)
		}
		if s := tag.Get("short"); s != "" {
			o = fmt.Sprintf("-%s, --%s%s", tag.Get("short"), tag.Get("long"), a)
		} else {
			o = fmt.Sprintf("--%s%s", tag.Get("long"), a)
		}

		desc := tag.Get("description")
		if i := strings.Index(desc, "\n"); i >= 0 {
			var buf bytes.Buffer
			buf.WriteString(desc[:i+1])
			desc = desc[i+1:]
			const indent = "                        "
			for {
				if i = strings.Index(desc, "\n"); i >= 0 {
					buf.WriteString(indent)
					buf.WriteString(desc[:i+1])
					desc = desc[i+1:]
					continue
				}
				break
			}
			if len(desc) > 0 {
				buf.WriteString(indent)
				buf.WriteString(desc)
			}
			desc = buf.String()
		}
		help = append(help, fmt.Sprintf("  %-15s %s", o, desc))
	}

	return help
}

func (c *CLI) showHelp() {
	opts := strings.Join(c.buildHelp([]string{
		"Major",
		"Minor",
		"Patch",
		"Pre",
		"Prefix",
		"Help",
		"Version",
	}), "\n")

	help := `
Usage: git-semv [--version] [--help] command <options>

Commands:
  now             Current version for semantic versionning
  next            Next version for semantic versionning
  bump            Bump version as next, and push to origin

Options:
%s
`
	fmt.Fprintf(c.outStream, help, opts)
}

func (c *CLI) target() string {
	if c.Major {
		return "major"
	}
	if c.Minor {
		return "minor"
	}
	if c.Patch {
		return "patch"
	}
	return ""
}

func (c *CLI) run(a []string) {
	p := flags.NewParser(c, flags.PrintErrors|flags.PassDoubleDash)
	args, err := p.ParseArgs(a)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return
	}

	if c.Help {
		c.showHelp()
		os.Exit(ExitErr)
		return
	}

	if c.Version {
		fmt.Fprintf(c.errStream, "git-semv version %s [%v, %v]\n", version, commit, date)
		os.Exit(ExitOK)
		return
	}

	c.Command = args[0]

	if len(args) > 1 {
		c.Args = args[1:]
	}

	switch c.Command {
	case "now":
		semv, err := semv.New(c.Prefix)
		if err != nil {
			fmt.Printf("%#v\n", err)
		}
		fmt.Printf("%s\n", semv)

	case "next":
		semv, err := semv.New(c.Prefix)
		if err != nil {
			fmt.Printf("%#v\n", err)
		}
		fmt.Printf("%s\n", semv.Bump(c.target(), c.Pre))

	case "bump":

	default:
		fmt.Fprintf(c.errStream, "Error: command is not available\n")
		c.showHelp()
		os.Exit(ExitErr)
		return
	}

	os.Exit(ExitOK)
}
