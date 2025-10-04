package semv

import (
	"os/exec"
)

// Cmder interface
type Cmder interface {
	Do(name string, arg ...string) ([]byte, error)
	DoWithEnv(name string, env []string, arg ...string) ([]byte, error)
}

// Cmd struct
type Cmd struct{}

// Do execute command
func (c Cmd) Do(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).Output()
}

// DoWithEnv execute command with environment variables
func (c Cmd) DoWithEnv(name string, env []string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	cmd.Env = env
	return cmd.Output()
}
