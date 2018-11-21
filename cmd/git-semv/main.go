package main

import (
	"fmt"

	gitsemv "github.com/linyows/git-semv"
)

func main() {
	v, _ := gitsemv.New()
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
