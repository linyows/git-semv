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

func TestBumpWhenNonPre(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	v := semv.Next("major", false)
	expected := "v3.0.0"

	if expected != v {
		t.Errorf("expected %s, but %s", expected, v)
	}
}

func TestBumpWhenPre(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	v := semv.Next("major", true)
	expected := "v3.0.0-rc.0"

	if expected != v {
		t.Errorf("expected %s, but %s", expected, v)
	}
}

func TestNextMajor(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	semv.nextMajor()
	expected := "3.0.0"

	if expected != semv.next.String() {
		t.Errorf("expected %s, but %s", expected, semv.next)
	}
}

func TestNextMinor(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	semv.nextMinor()
	expected := "2.4.0"

	if expected != semv.next.String() {
		t.Errorf("expected %s, but %s", expected, semv.next)
	}
}

func TestNextPatch(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	semv.nextPatch()
	expected := "2.3.5"

	if expected != semv.next.String() {
		t.Errorf("expected %s, but %s", expected, semv.next)
	}
}

func TestNextPreRelease(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New("v")
	semv.nextPreRelease()
	expected := "2.3.4-rc.3"

	if expected != semv.next.String() {
		t.Errorf("expected %s, but %s", expected, semv.next)
	}
}
