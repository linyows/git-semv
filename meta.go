package semv

import (
	"bytes"
)

// UsernameCmd for git
var UsernameCmd = []string{"config", "user.name"}

// LatestCommitCmd for git
var LatestCommitCmd = []string{"describe", "--always"}
var usernameCmder Cmder
var latestCommitCmder Cmder

func username() ([]byte, error) {
	if usernameCmder == nil {
		usernameCmder = Cmd{}
	}

	b, err := usernameCmder.Do(git, UsernameCmd...)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(b), nil
}

func latestCommit() ([]byte, error) {
	if latestCommitCmder == nil {
		latestCommitCmder = Cmd{}
	}

	b, err := latestCommitCmder.Do(git, LatestCommitCmd...)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(b), nil
}

func meta() ([]string, error) {
	user, err := username()
	if err != nil {
		return nil, err
	}

	hash, err := latestCommit()
	if err != nil {
		return nil, err
	}

	return []string{string(hash), string(user)}, nil
}
