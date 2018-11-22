package semv

import (
	"bytes"
	"strings"

	"github.com/blang/semver"
)

// TagCmd for tag list
var TagCmd = []string{"tag", "--list", "--sort=-v:refname"}
var git = "git"
var cmder Cmder
var defaultVersion = "0.0.0"
var defaultPreVersion = "0"
var defaultPreVersionPrefix = "rc"

// Semv struct
type Semv struct {
	Current semver.Version
	Next    semver.Version
}

// New returns Semv
func New() (*Semv, error) {
	if cmder == nil {
		cmder = Cmd{}
	}

	b, err := cmder.Do(git, TagCmd...)
	if err != nil {
		return nil, err
	}

	v := defaultVersion

	if len(bytes.TrimSpace(b)) > 0 {
		v = strings.Split(string(b), "\n")[0]
	}

	vv, err := semver.ParseTolerant(v)
	if err != nil {
		return nil, err
	}
	return &Semv{Current: vv}, err
}

// BumpMajor bump version for major version
func (v *Semv) BumpMajor() {
	copy := v.Current
	v.Next = copy
	v.Next.Major++
	v.Next.Minor = 0
	v.Next.Patch = 0
	v.Next.Pre = nil
}

// BumpMinor bump version for minor version
func (v *Semv) BumpMinor() {
	copy := v.Current
	v.Next = copy
	v.Next.Minor++
	v.Next.Patch = 0
	v.Next.Pre = nil
}

// BumpPatch bump version for patch version
func (v *Semv) BumpPatch() {
	copy := v.Current
	v.Next = copy
	v.Next.Patch++
	v.Next.Pre = nil
}

// BumpPre bump version for pre version
func (v *Semv) BumpPre() {
	copy := v.Current
	v.Next = copy

	if len(v.Next.Pre) > 0 {
		notB := true
		for i, pre := range v.Next.Pre {
			if pre.IsNumeric() {
				v.Next.Pre[i].VersionNum++
				notB = false
			}
		}
		if notB {
			p, err := semver.NewPRVersion(defaultPreVersion)
			if err == nil {
				v.Next.Pre = append(v.Next.Pre, p)
			}
		}
		return
	}

	prefix, err := semver.NewPRVersion(defaultPreVersionPrefix)
	if err == nil {
		v.Next.Pre = append(v.Next.Pre, prefix)
	}

	prever, err := semver.NewPRVersion(defaultPreVersion)
	if err == nil {
		v.Next.Pre = append(v.Next.Pre, prever)
	}
}
