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
	semv, err := New()
	if err != nil {
		t.Fatalf("Error: %#v", err)
	}
	expected := "12.345.67"
	if expected != semv.Current.String() {
		t.Errorf("expected %s, but %s", expected, semv.Current)
	}
}

func TestNewWhenTagNothing(t *testing.T) {
	cmder = MockedCmd{
		Out: `
`,
	}
	semv, err := New()
	if err != nil {
		t.Fatalf("Error: %#v", err)
	}
	expected := "0.0.0"
	if expected != semv.Current.String() {
		t.Errorf("expected %s, but %s", expected, semv.Current)
	}
}

func TestBumpMajor(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New()
	semv.BumpMajor()
	expected := "3.0.0"

	if expected != semv.Next.String() {
		t.Errorf("expected %s, but %s", expected, semv.Next)
	}
}

func TestBumpMinor(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New()
	semv.BumpMinor()
	expected := "2.4.0"

	if expected != semv.Next.String() {
		t.Errorf("expected %s, but %s", expected, semv.Next)
	}
}

func TestBumpPatch(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New()
	semv.BumpPatch()
	expected := "2.3.5"

	if expected != semv.Next.String() {
		t.Errorf("expected %s, but %s", expected, semv.Next)
	}
}

func TestBumpPre(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := New()
	semv.BumpPre()
	expected := "2.3.4-rc.3"

	if expected != semv.Next.String() {
		t.Errorf("expected %s, but %s", expected, semv.Next)
	}
}
