package main

import (
	"fmt"

	"github.com/linyows/git-semv/git"
)

func main() {
	v, _ := semv.New()
	fmt.Printf("%s\n", v.Current)
	v.BumpMajor()
	fmt.Printf("%s\n", v.Next)
	v.BumpMinor()
	fmt.Printf("%s\n", v.Next)
	v.BumpPatch()
	fmt.Printf("%s\n", v.Next)
	v.BumpPre()
	fmt.Printf("%s\n", v.Next)
}
