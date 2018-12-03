package semv

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
	unknownFlag := "Error: unknown flag `unknown'\n"
	unknownCmd := `Error: command is not available
`

	semvers := `v8.8.8
v12.0.1
v12.344.0+20130313144700
v12.345.66
v12.345.67
`

	semversWithPre := `v2.3.4-rc.2
v8.8.8
v12.0.1
v12.3.0-alpha
v12.3.0-alpha.0
v12.3.0-alpha.1
v12.3.0-alpha.1.beta
v12.3.0-beta
v12.3.0-beta.5
v12.3.0-rc
v12.344.0+20130313144700
v12.345.66
v12.345.67
v13.0.0-alpha.0
`

	tests := []struct {
		cmd   []string
		wantO []byte
		wantE []byte
		wantS int
	}{
		// list
		{[]string{}, []byte(semvers), []byte(""), ExitOK},
		{[]string{"list"}, []byte(semvers), []byte(""), ExitOK},
		{[]string{"-a"}, []byte(semversWithPre), []byte(""), ExitOK},
		{[]string{"--all"}, []byte(semversWithPre), []byte(""), ExitOK},
		{[]string{"list", "-a"}, []byte(semversWithPre), []byte(""), ExitOK},
		{[]string{"list", "--all"}, []byte(semversWithPre), []byte(""), ExitOK},
		// now
		{[]string{"now"}, []byte("v12.345.67\n"), []byte(""), ExitOK},
		{[]string{"latest"}, []byte("v12.345.67\n"), []byte(""), ExitOK},
		// next
		{[]string{"major"}, []byte("v13.0.0\n"), []byte(""), ExitOK},
		{[]string{"minor"}, []byte("v12.346.0\n"), []byte(""), ExitOK},
		{[]string{"patch"}, []byte("v12.345.68\n"), []byte(""), ExitOK},
		// pre
		{[]string{"major", "-p"}, []byte("v13.0.0-alpha.1\n"), []byte(""), ExitOK},
		{[]string{"major", "--pre"}, []byte("v13.0.0-alpha.1\n"), []byte(""), ExitOK},
		{[]string{"minor", "--pre"}, []byte("v12.346.0-alpha.0\n"), []byte(""), ExitOK},
		{[]string{"patch", "--pre"}, []byte("v12.345.68-alpha.0\n"), []byte(""), ExitOK},
		{[]string{"major", "--pre-name", "rc"}, []byte("v13.0.0-rc.0\n"), []byte(""), ExitOK},
		// build
		{[]string{"major", "-b"}, []byte("v13.0.0+2f994ff.foobar\n"), []byte(""), ExitOK},
		{[]string{"major", "--build"}, []byte("v13.0.0+2f994ff.foobar\n"), []byte(""), ExitOK},
		{[]string{"minor", "--build"}, []byte("v12.346.0+2f994ff.foobar\n"), []byte(""), ExitOK},
		{[]string{"patch", "--build"}, []byte("v12.345.68+2f994ff.foobar\n"), []byte(""), ExitOK},
		{[]string{"major", "--build-name", "baz"}, []byte("v13.0.0+baz\n"), []byte(""), ExitOK},
		// options
		{[]string{"-h"}, []byte(help), []byte(""), ExitErr},
		{[]string{"--help"}, []byte(help), []byte(""), ExitErr},
		{[]string{"-v"}, []byte(""), []byte(ver), ExitOK},
		{[]string{"--version"}, []byte(""), []byte(ver), ExitOK},
		// unknown
		{[]string{"--unknown=abc"}, []byte(""), []byte(unknownFlag), ExitErr},
		{[]string{"unknown"}, []byte(help), []byte(unknownCmd), ExitErr},
	}

	tagCmder = MockedCmd{Out: mixed}
	usernameCmder = MockedCmd{Out: "foobar\n", Err: ""}
	latestCommitCmder = MockedCmd{Out: "2f994ff\n", Err: ""}

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

func TestRunCLI(t *testing.T) {
	ver := "git-semv version dev [none, unknown]\n"
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	args := []string{"-v"}
	status := RunCLI(out, err, args)
	if status != 0 {
		t.Errorf("exit status = %d; want 0", status)
	}
	if bytes.Compare([]byte(""), out.Bytes()) != 0 {
		t.Errorf("output = %s; want empty", out)
	}
	if bytes.Compare([]byte(ver), err.Bytes()) != 0 {
		t.Errorf("err = %v; want %s", err, ver)
	}
}
