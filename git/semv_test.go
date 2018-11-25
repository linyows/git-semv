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

func TestCurrent(t *testing.T) {
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
	semv, err := Current()
	if err != nil {
		t.Fatalf("Error: %#v", err)
	}
	expected := "v12.345.67"
	if expected != semv.String() {
		t.Errorf("expected %s, but %s", expected, semv)
	}
}

func TestNewWhenTagNothing(t *testing.T) {
	cmder = MockedCmd{
		Out: `
`,
	}
	semv, err := Current()
	if err != nil {
		t.Fatalf("Error: %#v", err)
	}
	expected := "v0.0.0"
	if expected != semv.String() {
		t.Errorf("expected %s, but %s", expected, semv)
	}
}

func TestString(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, err := Current()
	if err != nil {
		t.Fatalf("Error: %#v", err)
	}
	expected := "v2.3.4-rc.2"

	if expected != semv.String() {
		t.Errorf("expected %s, but %s", expected, semv)
	}
}

func TestIncrementWhenNonPre(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := Current()
	semv.Next("major")
	expected := "v3.0.0"

	if expected != semv.String() {
		t.Errorf("expected %s, but %s", expected, semv)
	}
}

func TestIncrementWhenPre(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := Current()
	semv.Next("major").PreRelease("")
	expected := "v3.0.0-rc.0"

	if expected != semv.String() {
		t.Errorf("expected %s, but %s", expected, semv)
	}
}

func TestIncrementMajor(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := Current()
	semv.incrementMajor()
	expected := "v3.0.0"

	if expected != semv.String() {
		t.Errorf("expected %s, but %s", expected, semv)
	}
}

func TestIncrementMinor(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := Current()
	semv.incrementMinor()
	expected := "v2.4.0"

	if expected != semv.String() {
		t.Errorf("expected %s, but %s", expected, semv)
	}
}

func TestIncrementPatch(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := Current()
	semv.incrementPatch()
	expected := "v2.3.5"

	if expected != semv.String() {
		t.Errorf("expected %s, but %s", expected, semv)
	}
}

func TestPreRelease(t *testing.T) {
	cmder = MockedCmd{Out: `v2.3.4-rc.2`}
	semv, _ := Current()
	semv.PreRelease("")
	expected := "v2.3.4-rc.3"

	if expected != semv.String() {
		t.Errorf("expected %s, but %s", expected, semv)
	}
}
