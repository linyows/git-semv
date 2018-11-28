package semv

import (
	"bytes"
	"strings"

	"github.com/blang/semver"
)

// TagCmd for tag list
var TagCmd = []string{"tag", "--list", "--sort=v:refname"}
var git = "git"
var tagCmder Cmder
var defaultVersion = "0.0.0"

// List struct
type List struct {
	data semver.Versions
}

// GetList returns List
func GetList() (*List, error) {
	list, err := getVersions()
	if err != nil {
		return nil, err
	}
	return &List{data: list}, nil
}

// FindSimilar finds similar one
func (l *List) FindSimilar(v semver.Version) *Semv {
	for _, vv := range l.data {
		if vv.Major == v.Major && vv.Minor == v.Minor && vv.Patch == v.Patch {
			return MustNew(vv.String())
		}
	}
	return &Semv{}
}

// String to string
func (l *List) String() string {
	var ll []string
	for _, v := range l.data {
		ll = append(ll, defaultTagPrefix+v.String())
	}
	return strings.Join(ll, "\n")
}

// Latest gets latest version from List
func (l *List) Latest() *Semv {
	if l.data.Len() > 0 {
		return &Semv{data: l.data[len(l.data)-1]}
	}
	return &Semv{data: semver.MustParse(defaultVersion)}
}

// WithoutPreRelease excludes pre-release
func (l *List) WithoutPreRelease() *List {
	var list semver.Versions
	for _, v := range l.data {
		if len(v.Pre) > 0 {
			continue
		}
		list = append(list, v)
	}
	return &List{data: list}
}

// OnlyPreRelease includes only pre-release
func (l *List) OnlyPreRelease() *List {
	var list semver.Versions
	for _, v := range l.data {
		if len(v.Pre) == 0 {
			continue
		}
		list = append(list, v)
	}
	return &List{data: list}
}

// getVersions executes git tag as command
func getVersions() (semver.Versions, error) {
	if tagCmder == nil {
		tagCmder = Cmd{}
	}

	b, err := tagCmder.Do(git, TagCmd...)
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
		if defaultTagPrefix != "" {
			trimmed := strings.TrimPrefix(v, defaultTagPrefix)
			if trimmed == v {
				continue
			}
		}
		sv, err := semver.ParseTolerant(v)
		if err != nil {
			continue
		}
		list = append(list, sv)
	}
	semver.Sort(list)

	return list, nil
}
