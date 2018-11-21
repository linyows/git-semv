package gitsemv

import (
	"os/exec"
	"strings"

	"github.com/blang/semver"
)

type GitSemv struct {
	Current semver.Version
	Next    semver.Version
}

func New() (*GitSemv, error) {
	b, err := exec.Command("git", "tag", "--list", "--sort=-v:refname").Output()
	if err != nil {
		return nil, err
	}
	v := strings.Split(string(b), "\n")[0]
	vv, err := semver.ParseTolerant(v)
	if err != nil {
		return nil, err
	}
	return &GitSemv{Current: vv}, err
}

func (v *GitSemv) BumpMajor() {
	copy := v.Current
	v.Next = copy
	v.Next.Major += 1
	v.Next.Minor = 0
	v.Next.Patch = 0
	v.Next.Pre = nil
}

func (v *GitSemv) BumpMinor() {
	copy := v.Current
	v.Next = copy
	v.Next.Minor += 1
	v.Next.Patch = 0
	v.Next.Pre = nil
}

func (v *GitSemv) BumpPatch() {
	copy := v.Current
	v.Next = copy
	v.Next.Patch += 1
	v.Next.Pre = nil
}

func (v *GitSemv) BumpPre() {
	copy := v.Current
	v.Next = copy
	if len(v.Next.Pre) > 0 {
		notB := true
		for _, pre := range v.Next.Pre {
			if pre.IsNumeric() {
				pre.VersionNum += 1
				notB = false
			}
		}
		if notB {
			p, err := semver.NewPRVersion("0")
			if err == nil {
				v.Next.Pre = append(v.Next.Pre, p)
			}
		}
	} else {
		p, err := semver.NewPRVersion("rc.0")
		if err == nil {
			v.Next.Pre = append(v.Next.Pre, p)
		}
	}
}
