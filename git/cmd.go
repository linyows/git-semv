package semv

import (
	"os/exec"
)

// Cmder interface
type Cmder interface {
	Do(name string, arg ...string) ([]byte, error)
}

// Cmd struct
type Cmd struct{}

// Do execute command
func (c Cmd) Do(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).Output()
}
