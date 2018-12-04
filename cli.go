package semv

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

const (
	// ExitOK for exit code
	ExitOK int = 0

	// ExitErr for exit code
	ExitErr int = 1
)

// Env struct
type Env struct {
	Out, Err io.Writer
	Args     []string
	Version  string
	Commit   string
	Date     string
}

type cli struct {
	env       Env
	command   string
	Pre       bool   `long:"pre" short:"p" description:"Pre-Release version indicates(ex: 0.0.1-rc.0)"`
	PreName   string `long:"pre-name" description:"Specify pre-release version name"`
	Build     bool   `long:"build" short:"b" description:"Build version indicates(ex: 0.0.1+3222d31.foo)"`
	BuildName string `long:"build-name" description:"Specify build version name"`
	All       bool   `long:"all" short:"a" description:"Include everything such as pre-release and build versions in list"`
	Bump      bool   `long:"bump" short:"B" description:"Create tag and Push to origin"`
	Prefix    string `long:"prefix" short:"x" description:"Prefix for version and tag(default: v)"`
	Help      bool   `long:"help" short:"h" description:"Show this help message and exit"`
	Version   bool   `long:"version" short:"v" description:"Prints the version number"`
}

var gitTagCmder Cmder
var gitPushTagCmder Cmder

// RunCLI runs for CLI
func RunCLI(env Env) int {
	return (&cli{env: env}).run()
}

func (c *cli) buildHelp(names []string) []string {
	var help []string
	t := reflect.TypeOf(cli{})

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
		help = append(help, fmt.Sprintf("  %-18s %s", o, desc))
	}

	return help
}

func (c *cli) showHelp() {
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
	fmt.Fprintf(c.env.Out, help, opts)
}

func (c *cli) run() int {
	p := flags.NewParser(c, flags.PassDoubleDash)
	args, err := p.ParseArgs(c.env.Args)
	if err != nil {
		fmt.Fprintf(c.env.Err, "Error: %s\n", err)
		return ExitErr
	}

	if c.Help {
		c.showHelp()
		return ExitErr
	}

	if c.Version {
		fmt.Fprintf(c.env.Err, "git-semv version %s [%v, %v]\n", c.env.Version, c.env.Commit, c.env.Date)
		return ExitOK
	}

	if len(args) > 0 {
		c.command = args[0]
	} else {
		c.command = "list"
	}

	switch c.command {
	case "list":
		list, err := GetList()
		if err != nil {
			fmt.Fprintf(c.env.Err, "Error: %s\n", err)
		}
		if c.All == false {
			list = list.WithoutPreRelease()
		}
		fmt.Fprintf(c.env.Out, "%s\n", list)

	case "now", "latest":
		latest, err := Latest()
		if err != nil {
			fmt.Fprintf(c.env.Err, "Error: %s\n", err)
		}
		fmt.Fprintf(c.env.Out, "%s\n", latest)

	case "major", "minor", "patch":
		latest, err := Latest()
		if err != nil {
			fmt.Fprintf(c.env.Err, "Error: %s\n", err)
		}
		next := latest.Next(c.command)
		if c.Pre || c.PreName != "" {
			next.PreRelease(c.PreName)
		}
		if c.Build || c.BuildName != "" {
			next.Build(c.BuildName)
		}
		if c.Bump {
			if gitTagCmder == nil {
				gitTagCmder = Cmd{}
			}
			_, err = gitTagCmder.Do("git", "tag", next.String())
			if err != nil {
				fmt.Fprintf(c.env.Err, "Error: %s\n", err)
				return ExitErr
			}
			if gitPushTagCmder == nil {
				gitPushTagCmder = Cmd{}
			}
			_, err = gitPushTagCmder.Do("git", "push", "origin", next.String())
			if err != nil {
				fmt.Fprintf(c.env.Err, "Error: %s\n", err)
				return ExitErr
			}
			fmt.Fprintf(c.env.Out, "Bumped version to %s\n", next)
		} else {
			fmt.Fprintf(c.env.Out, "%s\n", next)
		}

	default:
		fmt.Fprintf(c.env.Err, "Error: command is not available: %s\n", c.command)
		c.showHelp()
		return ExitErr
	}

	return ExitOK
}
