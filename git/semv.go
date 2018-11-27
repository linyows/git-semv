package semv

import (
	"bytes"
	"fmt"

	"github.com/blang/semver"
)

// UsernameCmd for git
var UsernameCmd = []string{"config", "user.name"}

// LatestCommitCmd for git
var LatestCommitCmd = []string{"describe", "--always"}
var usernameCmder Cmder
var latestCommitCmder Cmder
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

// Latest returns latest version
func Latest() (*Semv, error) {
	list, err := GetList()
	if err != nil {
		return nil, err
	}
	return &Semv{
		data: list.WithoutPreRelease().Latest(),
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
	list := v.list.OnlyPreRelease()

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
func (v *Semv) Build(name string) (*Semv, error) {
	if name == "" {
		m, err := meta()
		if err != nil {
			return nil, err
		}
		v.data.Build = m
	} else {
		v.data.Build = []string{name}
	}

	return v, nil
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

func username() ([]byte, error) {
	if usernameCmder == nil {
		usernameCmder = Cmd{}
	}

	b, err := usernameCmder.Do(git, UsernameCmd...)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(b), nil
}

func latestCommit() ([]byte, error) {
	if latestCommitCmder == nil {
		latestCommitCmder = Cmd{}
	}

	b, err := latestCommitCmder.Do(git, LatestCommitCmd...)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(b), nil
}

func meta() ([]string, error) {
	user, err := username()
	if err != nil {
		return nil, err
	}

	hash, err := latestCommit()
	if err != nil {
		return nil, err
	}

	return []string{string(hash), string(user)}, nil
}
