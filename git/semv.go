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
	list    semver.Versions
}

// New returns Semv
func New(p string) (*Semv, error) {
	var sv semver.Version
	list, err := List()

	if len(list) > 0 {
		sv = list[len(list)-1]
	} else {
		sv, err = semver.ParseTolerant(defaultVersion)
		if err != nil {
			return nil, err
		}
	}

	copied := sv
	return &Semv{
		prefix:  p,
		current: sv,
		next:    copied,
		list:    list,
	}, err
}

// List returns list by git tag command
func List() (semver.Versions, error) {
	if cmder == nil {
		cmder = Cmd{}
	}

	b, err := cmder.Do(git, TagCmd...)
	if err != nil {
		return nil, err
	}
	b = bytes.TrimSpace(b)

	vv := []string{defaultVersion}
	if len(b) > 0 {
		vv = strings.Split(string(b), "\n")
	}

	var list semver.Versions
	for _, v := range vv {
		sv, err := semver.ParseTolerant(v)
		if err != nil {
			continue
		}
		list = append(list, sv)
	}
	semver.Sort(list)

	return list, nil
}

// List returns list
func (v *Semv) List(all bool) semver.Versions {
	if all == true {
		return v.list
	}

	var list semver.Versions
	for _, vv := range v.list {
		if len(vv.Pre) > 0 {
			continue
		}
		list = append(list, vv)
	}
	return list
}

// ListString returns list as string
func (v *Semv) ListString(all bool) string {
	var list []string
	for _, vv := range v.list {
		if all == false && len(vv.Pre) > 0 {
			continue
		}
		list = append(list, v.prefix+vv.String())
	}
	return strings.Join(list, "\n")
}

// String returns current version
func (v *Semv) String() string {
	return v.Current()
}

// Current returns current version
func (v *Semv) Current() string {
	return v.prefix + v.current.String()
}

// Next returns next version
func (v *Semv) Next() string {
	return v.prefix + v.next.String()
}

// Increment returns next version
func (v *Semv) Increment(target string) *Semv {
	switch target {
	case "major":
		v.incrementMajor()
	case "minor":
		v.incrementMinor()
	case "patch":
		v.incrementPatch()
	}
	return v
}

func (v *Semv) incrementMajor() {
	v.next.Major++
	v.next.Minor = 0
	v.next.Patch = 0
	v.next.Pre = nil
}

func (v *Semv) incrementMinor() {
	v.next.Minor++
	v.next.Patch = 0
	v.next.Pre = nil
}

func (v *Semv) incrementPatch() {
	v.next.Patch++
	v.next.Pre = nil
}

func (v *Semv) PreRelease(name string) {
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

func (v *Semv) Build(name string) {
}
