package semv

import (
	"fmt"
	"testing"
)

type MockedCmd struct {
	Out string
	Err string
}

func (c MockedCmd) Do(name string, arg ...string) ([]byte, error) {
	var err error
	if c.Err != "" {
		err = fmt.Errorf(c.Err)
	}
	return []byte(c.Out), err
}

var mixed = `v13.0.0-alpha.0
v12.345.67
v12.345.66
v12.344.0+20130313144700
v12.3.0-rc
v12.3.0-beta.5
v12.3.0-beta
v12.3.0-alpha.1.beta
v12.3.0-alpha.1
v12.3.0-alpha.0
v12.3.0-alpha
v12.0.1
v8.8.8
v2.3.4-rc.2
foo
bar-0
1.0.0
`

var empty = `
`

func TestCurrent(t *testing.T) {
	tests := []struct {
		out  string
		want string
	}{
		{mixed, "v12.345.67"},
		{empty, "v0.0.0"},
	}

	for i, tt := range tests {
		tagCmder = MockedCmd{Out: tt.out}
		v, err := Current()
		if err != nil {
			t.Fatal(err)
		}
		if v.String() != tt.want {
			t.Errorf("test[%d]: out = %s, Semv = %s; want %s", i, tt.out, v, tt.want)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		v    *Semv
		want string
	}{
		{&Semv{}, "v0.0.0"},
		{MustNew("1.0.0"), "v1.0.0"},
	}

	for i, tt := range tests {
		if tt.v.String() != tt.want {
			t.Errorf("test[%d]: Semv(%#v) = %s; want %s", i, tt.v, tt.v, tt.want)
		}
	}
}

func TestNext(t *testing.T) {
	tagCmder = MockedCmd{Out: mixed}
	v, err := Current()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		target string
		want   string
	}{
		{"major", "v13.0.0"},
		{"minor", "v12.346.0"},
		{"patch", "v12.345.68"},
	}

	for i, tt := range tests {
		vn := v.Next(tt.target)
		if vn.String() != tt.want {
			t.Errorf("test[%d]: Semv(%#v) = %s; want %s", i, vn, vn, tt.want)
		}
	}
}

func TestPreRelease(t *testing.T) {
	tagCmder = MockedCmd{Out: mixed}
	v, err := Current()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		target  string
		preName string
		want    string
	}{
		{"major", "", "v13.0.0-alpha.1"},
		{"major", "beta", "v13.0.0-beta.0"},
		{"minor", "", "v12.346.0-alpha.0"},
		{"minor", "beta", "v12.346.0-beta.0"},
		{"patch", "", "v12.345.68-alpha.0"},
		{"patch", "beta", "v12.345.68-beta.0"},
	}

	for i, tt := range tests {
		vn, err := v.Next(tt.target).PreRelease(tt.preName)
		if err != nil {
			t.Fatal(err)
		}
		if vn.String() != tt.want {
			t.Errorf("test[%d]: target = %s; name = %s; Semv = %s; want %s", i, tt.target, tt.preName, vn, tt.want)
		}
	}
}

func TestIncrementMajor(t *testing.T) {
	tests := []struct {
		v    *Semv
		want string
	}{
		{MustNew("1.0.0"), "v2.0.0"},
		{MustNew("1.2.3"), "v2.0.0"},
		{MustNew("0.0.0"), "v1.0.0"},
	}

	for i, tt := range tests {
		tt.v.incrementMajor()
		if tt.v.String() != tt.want {
			t.Errorf("test[%d]: Semv(%#v) = %s; want %s", i, tt.v, tt.v, tt.want)
		}
	}
}

func TestIncrementMinor(t *testing.T) {
	tests := []struct {
		v    *Semv
		want string
	}{
		{MustNew("1.0.0"), "v1.1.0"},
		{MustNew("1.2.3"), "v1.3.0"},
		{MustNew("0.0.0"), "v0.1.0"},
	}

	for i, tt := range tests {
		tt.v.incrementMinor()
		if tt.v.String() != tt.want {
			t.Errorf("test[%d]: Semv(%#v) = %s; want %s", i, tt.v, tt.v, tt.want)
		}
	}
}

func TestIncrementPatch(t *testing.T) {
	tests := []struct {
		v    *Semv
		want string
	}{
		{MustNew("1.0.0"), "v1.0.1"},
		{MustNew("1.2.3"), "v1.2.4"},
		{MustNew("0.0.0"), "v0.0.1"},
	}

	for i, tt := range tests {
		tt.v.incrementPatch()
		if tt.v.String() != tt.want {
			t.Errorf("test[%d]: Semv(%#v) = %s; want %s", i, tt.v, tt.v, tt.want)
		}
	}
}
