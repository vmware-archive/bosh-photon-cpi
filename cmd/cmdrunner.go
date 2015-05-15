package cmd

import (
	"os/exec"
)

type Runner interface {
	Run(name string, args ...string) ([]byte, error)
}

type defaultRunner struct {
}

func NewRunner() Runner {
	return &defaultRunner{}
}

func (_ defaultRunner) Run(name string, args ...string) (out []byte, err error) {
	cmd := exec.Command(name, args...)
	out, err = cmd.CombinedOutput()
	return
}
