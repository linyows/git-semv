package semv

import (
	"testing"

	"github.com/blang/semver"
)

func TestGetList(t *testing.T) {
	want := `v2.3.4-rc.2
v8.8.8
v12.0.1
v12.3.0-alpha
v12.3.0-alpha.0
v12.3.0-alpha.1
v12.3.0-alpha.1.beta
v12.3.0-beta
v12.3.0-beta.5
v12.3.0-rc
v12.344.0+20130313144700
v12.345.66
v12.345.67
v13.0.0-alpha.0`

	tests := []struct {
		out  string
		want string
	}{
		{mixed, want},
		{empty, ""},
	}

	for i, tt := range tests {
		tagCmder = MockedCmd{Out: tt.out}
		l, err := GetList()
		if err != nil {
			t.Fatal(err)
		}
		if l.String() != tt.want {
			t.Errorf("test[%d]: List = %s; want %s", i, l, tt.want)
		}
	}
}

func TestFindSimilar(t *testing.T) {
	out := `v1.0.0
v1.0.1
v1.0.2-rc.0`

	tests := []struct {
		out  string
		find string
		want string
	}{
		{out, "v1.0.2", "v1.0.2-rc.0"},
		{out, "v1.0.3", ""},
		{"", "v1.0.3", ""},
	}

	for i, tt := range tests {
		tagCmder = MockedCmd{Out: tt.out}
		l, err := GetList()
		if err != nil {
			t.Fatal(err)
		}
		v, err := semver.ParseTolerant(tt.find)
		if err != nil {
			t.Fatal(err)
		}
		semv := l.FindSimilar(v)
		if semv.String() != tt.want {
			t.Errorf("test[%d]: semver.Version = %s; want %s", i, semv, tt.want)
		}
	}
}
