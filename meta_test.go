package semv

import (
	"reflect"
	"testing"
)

func TestUsername(t *testing.T) {
	w := "foobar"
	usernameCmder = MockedCmd{Out: w + "\n"}
	u, err := username()
	if err != nil {
		t.Fatalf("Error: %#v", err)
	}
	if w != string(u) {
		t.Errorf("username = %s, want %s", u, w)
	}
}

func TestLatestCommit(t *testing.T) {
	w := "2f994ff"
	latestCommitCmder = MockedCmd{Out: w + "\n"}
	h, err := latestCommit()
	if err != nil {
		t.Fatalf("Error: %#v", err)
	}
	if w != string(h) {
		t.Errorf("hash = %s, want %s", h, w)
	}
}

func TestMeta(t *testing.T) {
	u := "foobar"
	h := "2f994ff"
	usernameCmder = MockedCmd{Out: u + "\n"}
	latestCommitCmder = MockedCmd{Out: h + "\n"}
	m, err := meta()
	if err != nil {
		t.Fatalf("Error: %#v", err)
	}
	w := []string{h, u}
	if !reflect.DeepEqual(w, m) {
		t.Errorf("meta = %s, want %s", m, w)
	}
}
