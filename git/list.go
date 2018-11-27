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
func (l *List) FindSimilar(v semver.Version) semver.Version {
	for _, vv := range l.data {
		if vv.Major == v.Major && vv.Minor == v.Minor && vv.Patch == v.Patch {
			return semver.MustParse(vv.String())
		}
	}
	return semver.Version{}
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
func (l *List) Latest() semver.Version {
	if len(l.data) > 0 {
		return l.data[len(l.data)-1]
	}
	v, err := semver.Parse(defaultVersion)
	if err != nil {
		panic(err)
	}
	return v
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
		sv, err := semver.ParseTolerant(v)
		if err != nil {
			continue
		}
		list = append(list, sv)
	}
	semver.Sort(list)

	return list, nil
}
