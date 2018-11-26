package semv

import (
	"fmt"

	"github.com/blang/semver"
)

var defaultPreVersion = "0"
var defaultPreVersionPrefix = "alpha"
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
	vv := *v
	copied := &vv
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
func (v *Semv) PreRelease(name string) (*Semv, error) {
	list, err := NewPreReleaseList()
	if err != nil {
		return nil, err
	}

	prefix := name
	if name == "" {
		prefix = defaultPreVersionPrefix
	}

	same := list.FindSame(v.data)
	if same.String() != "0.0.0" {
		v.data = same
	}

	if len(v.data.Pre) > 0 {
		incremented := false
		mustIncremnt := false

		for i, pre := range v.data.Pre {
			if pre.IsNumeric() && mustIncremnt && i < 3 {
				v.data.Pre[i].VersionNum++
				incremented = true
			} else if pre.IsNumeric() == false && i == 0 {
				if pre.VersionStr == prefix {
					mustIncremnt = true
				} else if pre.Compare(semver.PRVersion{VersionStr: prefix, IsNum: false}) == 1 {
					return nil, fmt.Errorf("%s is less than %s", prefix, pre.VersionStr)
				} else {
					v.data.Pre[i].VersionStr = prefix
					incremented = true
				}
			}
		}

		if incremented == true {
			return v, nil
		}
	}

	prV, err := semver.NewPRVersion(prefix)
	if err == nil {
		v.data.Pre = append(v.data.Pre, prV)
	}

	prever, err := semver.NewPRVersion(defaultPreVersion)
	if err == nil {
		v.data.Pre = append(v.data.Pre, prever)
	}

	return v, nil
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
