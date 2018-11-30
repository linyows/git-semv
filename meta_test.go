package semv

import (
	"reflect"
	"testing"
)

func TestUsername(t *testing.T) {
	tests := []struct {
		out  string
		err  string
		want string
	}{
		{"foobar\n", "", "foobar"},
		{"", "failed", ""},
	}

	for i, tt := range tests {
		usernameCmder = MockedCmd{Out: tt.out, Err: tt.err}
		u, err := username()

		if (tt.err == "" && err != nil) || (tt.err != "" && err.Error() != tt.err) {
			t.Errorf("test[%d]: err = %s; want %s", i, err, tt.err)
		}

		if tt.want != string(u) {
			t.Errorf("test[%d]: username = %s; want %s", i, u, tt.want)
		}
	}
}

func TestLatestCommit(t *testing.T) {
	tests := []struct {
		out  string
		err  string
		want string
	}{
		{"2f994ff\n", "", "2f994ff"},
		{"", "failed", ""},
	}

	for i, tt := range tests {
		latestCommitCmder = MockedCmd{Out: tt.out, Err: tt.err}
		h, err := latestCommit()
		if (tt.err == "" && err != nil) || (tt.err != "" && err.Error() != tt.err) {
			t.Errorf("test[%d]: err = %s; want %s", i, err, tt.err)
		}

		if tt.want != string(h) {
			t.Errorf("test[%d]: hash = %s; want %s", i, h, tt.want)
		}
	}
}

func TestMeta(t *testing.T) {
	tests := []struct {
		out1 string
		err1 string
		out2 string
		err2 string
		want []string
	}{
		{"foobar\n", "", "2f994ff\n", "", []string{"2f994ff", "foobar"}},
		{"foobar\n", "", "", "failed", nil},
		{"", "failed", "2f994ff\n", "", nil},
	}
	for i, tt := range tests {
		usernameCmder = MockedCmd{Out: tt.out1, Err: tt.err1}
		latestCommitCmder = MockedCmd{Out: tt.out2, Err: tt.err2}
		m, err := meta()
		if (tt.err1 != "" || tt.err2 != "") && err == nil {
			t.Errorf("test[%d]: err = nil; want %s", i, tt.err1+tt.err2)
		}
		if tt.err1 == "" && tt.err2 == "" && err != nil {
			t.Errorf("test[%d]: err = %s; want no error", i, err)
		}
		if !reflect.DeepEqual(tt.want, m) {
			t.Errorf("test[%d]: meta = %s; want %s", i, m, tt.want)
		}
	}
}
