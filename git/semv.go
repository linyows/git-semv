package semv

import (
	"github.com/blang/semver"
)

var defaultPreVersion = "0"
var defaultPreVersionPrefix = "rc"
var defaultTagPrefix = "v"

// Semv struct
type Semv struct {
	data           semver.Version
	list           *List
	preReleaseName string
	buildName      string
}

// MustNew creates Semv
func MustNew(s string) *Semv {
	v, err := semver.ParseTolerant(s)
	if err != nil {
		panic(err)
	}
	return &Semv{data: v}
}

// Current returns current version
func Current() (*Semv, error) {
	list, err := NewStrictList()
	if err != nil {
		return nil, err
	}
	return &Semv{
		data: list.Current(),
		list: list,
	}, nil
}

// String to string
func (v *Semv) String() string {
	return defaultTagPrefix + v.data.String()
}

// Next returns next version
func (v *Semv) Next(target string) *Semv {
	copied := v
	switch target {
	case "major":
		copied.incrementMajor()
	case "minor":
		copied.incrementMinor()
	case "patch":
		copied.incrementPatch()
	}
	return copied
}

// PreRelease retuns
func (v *Semv) PreRelease(name string) {
	if len(v.data.Pre) > 0 {
		notB := true
		for i, pre := range v.data.Pre {
			if pre.IsNumeric() {
				v.data.Pre[i].VersionNum++
				notB = false
			}
		}
		if notB {
			p, err := semver.NewPRVersion(defaultPreVersion)
			if err == nil {
				v.data.Pre = append(v.data.Pre, p)
			}
		}
		return
	}

	prefix, err := semver.NewPRVersion(defaultPreVersionPrefix)
	if err == nil {
		v.data.Pre = append(v.data.Pre, prefix)
	}

	prever, err := semver.NewPRVersion(defaultPreVersion)
	if err == nil {
		v.data.Pre = append(v.data.Pre, prever)
	}
}

// Build retuns
func (v *Semv) Build(name string) {
}

func (v *Semv) incrementMajor() {
	v.data.Major++
	v.data.Minor = 0
	v.data.Patch = 0
	v.data.Pre = nil
}

func (v *Semv) incrementMinor() {
	v.data.Minor++
	v.data.Patch = 0
	v.data.Pre = nil
}

func (v *Semv) incrementPatch() {
	v.data.Patch++
	v.data.Pre = nil
}
