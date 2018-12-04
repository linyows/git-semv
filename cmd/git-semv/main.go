package main

import (
	"os"

	semv "github.com/linyows/git-semv"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	os.Exit(semv.RunCLI(semv.Env{
		Out:     os.Stdout,
		Err:     os.Stderr,
		Args:    os.Args,
		Version: version,
		Commit:  commit,
		Date:    date,
	}))
}
