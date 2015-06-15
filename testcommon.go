package main

import (
	"strings"
)

// Mock cmd.Runner implementation. The cmds map matches command prefixes to
// command output. E.g. if you want to mock out a call to "ls" and "file",
// pass in something like map[string]string{"ls": "stdout for ls", "file": "stdout for file"}
// When Run is called it will return output for the first key that is a substring of the
// name argument.
type fakeRunner struct {
	cmds map[string]string
}

func (r *fakeRunner) Run(name string, args ...string) (out []byte, err error) {
	for key := range r.cmds {
		if strings.HasPrefix(name, key) {
			return []byte(r.cmds[key]), nil
		}
	}
	return
}
