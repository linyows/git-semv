package main

import (
	"bytes"
	"fmt"
	"io"
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
	Pre                  bool   `long:"pre" short:"p" description:"Pre-Release version indicates(ex: 0.0.1-rc.0)"`
	PreName              string `long:"pre-name" description:"Specify pre-release version name"`
	Build                bool   `long:"build" short:"b" description:"Build version indicates(ex: 0.0.1+3222d31.foo)"`
	BuildName            string `long:"build-name" description:"Specify build version name"`
	All                  bool   `long:"all" short:"a" description:"Include everything such as pre-release and build versions in list"`
	Bump                 bool   `long:"bump" short:"B" description:"Create tag and Push to origin"`
	Prefix               string `long:"prefix" short:"x" description:"Prefix for version and tag(default: v)"`
	Help                 bool   `long:"help" short:"h" description:"Show this help message and exit"`
	Version              bool   `long:"version" short:"v" description:"Prints the version number"`
}

var BumpCmd = "git tag %s && git push origin %s"

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
			o = fmt.Sprintf("    --%s%s", tag.Get("long"), a)
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
		help = append(help, fmt.Sprintf("  %-18s %s", o, desc))
	}

	return help
}

func (c *CLI) showHelp() {
	opts := strings.Join(c.buildHelp([]string{
		"Pre",
		"PreRelease",
		"Build",
		"BuildName",
		"All",
		"Bump",
		"Prefix",
		"Help",
		"Version",
	}), "\n")

	help := `
Usage: git-semv [--version] [--help] command <options>

Commands:
  list               Sorted versions
  now, latest        Latest version
  major              Next major version: vX.0.0
  minor              Next minor version: v0.X.0
  patch              Next patch version: v0.0.X

Options:
%s
`
	fmt.Fprintf(c.outStream, help, opts)
}

func (c *CLI) run(a []string) int {
	p := flags.NewParser(c, flags.PassDoubleDash)
	args, err := p.ParseArgs(a)
	if err != nil {
		fmt.Fprintf(c.errStream, "Error: %s\n", err)
		return ExitErr
	}

	if c.Help {
		c.showHelp()
		return ExitErr
	}

	if c.Version {
		fmt.Fprintf(c.errStream, "git-semv version %s [%v, %v]\n", version, commit, date)
		return ExitOK
	}

	if len(args) > 0 {
		c.Command = args[0]
	} else {
		c.Command = "list"
	}

	if len(args) > 1 {
		c.Args = args[1:]
	}

	switch c.Command {
	case "list":
		list, err := semv.GetList()
		if err != nil {
			fmt.Fprintf(c.errStream, "Error: %s\n", err)
		}
		if c.All == false {
			list = list.WithoutPreRelease()
		}
		fmt.Fprintf(c.outStream, "%s\n", list)

	case "now", "latest":
		latest, err := semv.Latest()
		if err != nil {
			fmt.Fprintf(c.errStream, "Error: %s\n", err)
		}
		fmt.Fprintf(c.outStream, "%s\n", latest)

	case "major", "minor", "patch":
		latest, err := semv.Latest()
		if err != nil {
			fmt.Fprintf(c.errStream, "Error: %s\n", err)
		}
		next := latest.Next(c.Command)
		if c.Pre {
			next.PreRelease(c.PreName)
		}
		if c.Build {
			next.Build(c.BuildName)
		}
		if c.Bump {
			fmt.Fprintf(c.outStream, BumpCmd, next, next)
		} else {
			fmt.Fprintf(c.outStream, "%s\n", next)
		}

	default:
		fmt.Fprintf(c.errStream, "Error: command is not available\n")
		c.showHelp()
		return ExitErr
	}

	return ExitOK
}
