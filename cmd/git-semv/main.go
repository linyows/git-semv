package main

import (
	"os"

	semv "github.com/linyows/git-semv"
)

func main() {
	os.Exit(semv.RunCLI(os.Stdout, os.Stderr, os.Args[1:]))
}
