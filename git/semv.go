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
	prefix  string
	current semver.Version
	next    semver.Version
}

// New returns Semv
func New(p string) (*Semv, error) {
	if cmder == nil {
		cmder = Cmd{}
	}

	b, err := cmder.Do(git, TagCmd...)
	if err != nil {
		return nil, err
	}

	strV := defaultVersion

	if len(bytes.TrimSpace(b)) > 0 {
		strV = strings.Split(string(b), "\n")[0]
	}

	semV, err := semver.ParseTolerant(strV)
	if err != nil {
		return nil, err
	}

	copiedV := semV
	return &Semv{prefix: p, current: semV, next: copiedV}, err
}

// String return current version
func (v *Semv) String() string {
	return v.Current()
}

// Current return current version
func (v *Semv) Current() string {
	return v.prefix + v.current.String()
}

// Next return next version
func (v *Semv) Next() string {
	return v.prefix + v.next.String()
}

// Bump bump version by argument
func (v *Semv) Bump(target string, pre bool) string {
	switch target {
	case "major":
		v.BumpMajor()
	case "minor":
		v.BumpMinor()
	case "patch":
		v.BumpPatch()
	}
	if pre {
		v.BumpPre()
	}
	return v.Next()
}

// BumpMajor bump version for major version
func (v *Semv) BumpMajor() {
	v.next.Major++
	v.next.Minor = 0
	v.next.Patch = 0
	v.next.Pre = nil
}

// BumpMinor bump version for minor version
func (v *Semv) BumpMinor() {
	v.next.Minor++
	v.next.Patch = 0
	v.next.Pre = nil
}

// BumpPatch bump version for patch version
func (v *Semv) BumpPatch() {
	v.next.Patch++
	v.next.Pre = nil
}

// BumpPre bump version for pre version
func (v *Semv) BumpPre() {
	if len(v.next.Pre) > 0 {
		notB := true
		for i, pre := range v.next.Pre {
			if pre.IsNumeric() {
				v.next.Pre[i].VersionNum++
				notB = false
			}
		}
		if notB {
			p, err := semver.NewPRVersion(defaultPreVersion)
			if err == nil {
				v.next.Pre = append(v.next.Pre, p)
			}
		}
		return
	}

	prefix, err := semver.NewPRVersion(defaultPreVersionPrefix)
	if err == nil {
		v.next.Pre = append(v.next.Pre, prefix)
	}

	prever, err := semver.NewPRVersion(defaultPreVersion)
	if err == nil {
		v.next.Pre = append(v.next.Pre, prever)
	}
}
