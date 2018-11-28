package main

import (
	"bytes"
	"testing"
)

func TestCLIRun(t *testing.T) {
	ver := `git-semv version dev [none, unknown]
`
	help := `
Usage: git-semv [--version] [--help] command <options>

Commands:
  list               Sorted versions
  now, latest        Latest version
  major              Next major version: vX.0.0
  minor              Next minor version: v0.X.0
  patch              Next patch version: v0.0.X

Options:
  -p, --pre          Pre-Release version indicates(ex: 0.0.1-rc.0)
  -b, --build        Build version indicates(ex: 0.0.1+3222d31.foo)
      --build-name   Specify build version name
  -a, --all          Include everything such as pre-release and build versions in list
  -B, --bump         Create tag and Push to origin
  -x, --prefix       Prefix for version and tag(default: v)
  -h, --help         Show this help message and exit
  -v, --version      Prints the version number
`
	unknown := `Error: command is not available
`

	tests := []struct {
		cmd   []string
		wantO []byte
		wantE []byte
		wantS int
	}{
		{[]string{"-h"}, []byte(help), []byte(""), ExitErr},
		{[]string{"--help"}, []byte(help), []byte(""), ExitErr},
		{[]string{"-v"}, []byte(""), []byte(ver), ExitOK},
		{[]string{"--version"}, []byte(""), []byte(ver), ExitOK},
		{[]string{"unknown"}, []byte(help), []byte(unknown), ExitErr},
	}

	for i, tt := range tests {
		out, err := new(bytes.Buffer), new(bytes.Buffer)
		cli := &CLI{outStream: out, errStream: err}
		status := cli.run(tt.cmd)

		if status != tt.wantS {
			t.Errorf("test[%d]: status = %d; want %d", i, status, tt.wantS)
		}

		if bytes.Compare(tt.wantO, out.Bytes()) != 0 {
			t.Errorf("test[%d]: stdout = %s; want %s", i, out, tt.wantO)
		}

		if bytes.Compare(tt.wantE, err.Bytes()) != 0 {
			t.Errorf("test[%d]: stderr = %s; want %s", i, err, tt.wantE)
		}
	}
}
