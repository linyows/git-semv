package semv

import (
	"bytes"
	"strings"

	"github.com/blang/semver"
)

// TagCmd for tag list
var TagCmd = []string{"tag", "--list", "--sort=v:refname", "--format=%(refname:short)\t%(authordate:iso)\t%(subject)\t%(authorname)"}
var git = "git"
var tagCmder Cmder
var defaultVersion = "0.0.0"

// TagInfo holds tag information including metadata
type TagInfo struct {
	Version    semver.Version
	AuthorDate string
	Subject    string
	AuthorName string
}

// Format returns formatted tag info string
func (t *TagInfo) Format(tagName string) string {
	return tagName + "\t" + t.AuthorDate + "\t" + t.Subject + "\t" + t.AuthorName
}

// List struct
type List struct {
	data     semver.Versions
	tagInfos map[string]*TagInfo
}

// GetList returns List
func GetList() (*List, error) {
	list, tagInfos, err := getVersions()
	if err != nil {
		return nil, err
	}
	return &List{data: list, tagInfos: tagInfos}, nil
}

// FindSimilar finds similar one
func (l *List) FindSimilar(v semver.Version) *Semv {
	newSemv := &Semv{}
	for _, vv := range l.data {
		if vv.Major == v.Major && vv.Minor == v.Minor && vv.Patch == v.Patch {
			newSemv = MustNew(vv.String())
		}
	}
	return newSemv
}

// String to string
func (l *List) String() string {
	var ll []string
	for _, v := range l.data {
		tagName := defaultTagPrefix + v.String()
		if info, ok := l.tagInfos[tagName]; ok {
			ll = append(ll, info.Format(tagName))
		} else {
			ll = append(ll, tagName)
		}
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
	return &List{data: list, tagInfos: l.tagInfos}
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
	return &List{data: list, tagInfos: l.tagInfos}
}

// getVersions executes git tag as command
func getVersions() (semver.Versions, map[string]*TagInfo, error) {
	if tagCmder == nil {
		tagCmder = Cmd{}
	}

	b, err := tagCmder.Do(git, TagCmd...)
	if err != nil {
		return nil, nil, err
	}
	b = bytes.TrimSpace(b)

	vv := []string{defaultVersion}
	if len(b) > 0 {
		vv = strings.Split(string(b), "\n")
	}

	var list semver.Versions
	tagInfos := make(map[string]*TagInfo)

	for _, line := range vv {
		parts := strings.Split(line, "\t")
		tagName := parts[0]

		if defaultTagPrefix != "" {
			trimmed := strings.TrimPrefix(tagName, defaultTagPrefix)
			if trimmed == tagName {
				continue
			}
		}

		sv, err := semver.ParseTolerant(tagName)
		if err != nil {
			continue
		}

		list = append(list, sv)

		// Store tag info if we have author date, subject, and name
		if len(parts) >= 4 {
			authorDate := parts[1]
			if authorDate == "" {
				authorDate = "-"
			}
			subject := parts[2]
			if subject == "" {
				subject = "-"
			}
			authorName := parts[3]
			if authorName == "" {
				authorName = "-"
			}
			tagInfos[tagName] = &TagInfo{
				Version:    sv,
				AuthorDate: authorDate,
				Subject:    subject,
				AuthorName: authorName,
			}
		}
	}
	semver.Sort(list)

	return list, tagInfos, nil
}
