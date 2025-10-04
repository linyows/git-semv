package semv

import (
	"testing"
)

func TestCmdDo(t *testing.T) {
	cmd := Cmd{}
	out, err := cmd.Do("echo", "-n", "hi")
	if err != nil {
		t.Fatalf("Error: %#v", err)
	}
	expected := "hi"
	if expected != string(out) {
		t.Errorf("expected %s, but %s", expected, out)
	}
}

func TestCmdDoWithEnv(t *testing.T) {
	cmd := Cmd{}
	env := []string{"TEST_VAR=hello"}
	out, err := cmd.DoWithEnv("sh", env, "-c", "printf $TEST_VAR")
	if err != nil {
		t.Fatalf("Error: %#v", err)
	}
	expected := "hello"
	if expected != string(out) {
		t.Errorf("expected %s, but %s", expected, out)
	}
}
