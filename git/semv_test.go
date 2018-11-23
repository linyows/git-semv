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

func TestNew(t *testing.T) {
	cmder = MockedCmd{
		Out: `v12.345.67
v12.345.66
v12.344.0
v12.0.1
v8.8.8
foo
bar-0
1.0.0
`,
	}
	semv, err := New("v")
	if err != nil {
		t.Fatalf("Error: %#v", err)
	}
	expected := "v12.345.67"
	if expected != semv.Current() {
		t.Errorf("expected %s, but %s", expected, semv.Current())
	}
}

func TestNewWhenTagNothing(t *testing.T) {
	cmder = MockedCmd{
		Out: `
`,
	}
	semv, err := New("v")
	if err != nil {
		t.Fatalf("Error: %#v", err)
	}
	expected := "v0.0.0"
	if expected != semv.String() {
		t.Errorf("expected %s, but %s", expected, semv.String())
	}
}

func TestString(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	v := semv.String()
	expected := "v2.3.4-rc.2"

	if expected != v {
		t.Errorf("expected %s, but %s", expected, v)
	}
}

func TestCurrent(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	v := semv.Current()
	expected := "v2.3.4-rc.2"

	if expected != v {
		t.Errorf("expected %s, but %s", expected, v)
	}
}

func TestNext(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	v := semv.Next()
	expected := "v2.3.4-rc.2"

	if expected != v {
		t.Errorf("expected %s, but %s", expected, v)
	}
}

func TestBumpWhenNonPre(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	semv.Bump("major", false)
	expected := "v3.0.0"

	if expected != semv.Next() {
		t.Errorf("expected %s, but %s", expected, semv.Next())
	}
}

func TestBumpWhenPre(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	semv.Bump("major", true)
	expected := "v3.0.0-rc.0"

	if expected != semv.Next() {
		t.Errorf("expected %s, but %s", expected, semv.Next())
	}
}

func TestBumpMajor(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	semv.BumpMajor()
	expected := "v3.0.0"

	if expected != semv.Next() {
		t.Errorf("expected %s, but %s", expected, semv.Next())
	}
}

func TestBumpMinor(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	semv.BumpMinor()
	expected := "v2.4.0"

	if expected != semv.Next() {
		t.Errorf("expected %s, but %s", expected, semv.Next())
	}
}

func TestBumpPatch(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	semv.BumpPatch()
	expected := "v2.3.5"

	if expected != semv.Next() {
		t.Errorf("expected %s, but %s", expected, semv.Next())
	}
}

func TestBumpPre(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	semv.BumpPre()
	expected := "v2.3.4-rc.3"

	if expected != semv.Next() {
		t.Errorf("expected %s, but %s", expected, semv.Next())
	}
}
