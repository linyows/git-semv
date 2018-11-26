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

// List struct
type List struct {
	data semver.Versions
}

// strict for callback
type strict func(semver.Version) bool

// NewList returns List
func NewList() (*List, error) {
	list, err := fetch(func(v semver.Version) bool {
		return false
	})
	if err != nil {
		return nil, err
	}
	return &List{data: list}, nil
}

// NewStrictList returns strict List
func NewStrictList() (*List, error) {
	list, err := fetch(func(v semver.Version) bool {
		return len(v.Pre) > 0
	})
	if err != nil {
		return nil, err
	}
	return &List{data: list}, nil
}

// NewPreReleaseList returns only pre-release
func NewPreReleaseList() (*List, error) {
	list, err := fetch(func(v semver.Version) bool {
		return len(v.Pre) < 0
	})
	if err != nil {
		return nil, err
	}
	return &List{data: list}, nil
}

// FindSame finds same one
func (l *List) FindSame(v semver.Version) semver.Version {
	for _, vv := range l.data {
		if vv.Major == v.Major && vv.Minor == v.Minor && vv.Patch == v.Patch {
			return vv
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

// Current get current version from List
func (l *List) Current() semver.Version {
	if len(l.data) > 0 {
		return l.data[len(l.data)-1]
	}
	v, err := semver.ParseTolerant(defaultVersion)
	if err != nil {
		panic(err)
	}
	return v
}

// fetch fetches versions by git command
func fetch(fn strict) (semver.Versions, error) {
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
		if fn(sv) {
			continue
		}
		list = append(list, sv)
	}
	semver.Sort(list)

	return list, nil
}
