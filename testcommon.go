package main

import (
	"fmt"
	"github.com/esxcloud/bosh-esxcloud-cpi/logger"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega/types"
	"io/ioutil"
	"path/filepath"
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

func newLogger(desc ginkgo.GinkgoTestDescription) logger.Logger {
	// Prefix the log name with test filename and line number to make it easy to read test logs
	name := strings.TrimSuffix(filepath.Base(desc.FileName), filepath.Ext(desc.FileName))
	logger, _ := logger.New(fmt.Sprintf("%s_lineNum_%d_", name, desc.LineNumber))
	return logger
}

func containLogData() types.GomegaMatcher {
	return &logMatcher{}
}

type logMatcher struct {
}

func (m *logMatcher) Match(actual interface{}) (success bool, err error) {
	logFile, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("containLogData matcher expects a string indicating log path")
	}
	data, err := ioutil.ReadFile(logFile)
	if err != nil {
		return false, err
	}
	logData := string(data[:])
	if len(logData) < 1 {
		return false, fmt.Errorf("No log data found in log file: %s", logFile)
	}
	return true, nil
}

func (m *logMatcher) FailureMessage(actual interface{}) string {
	return "Expected log data to be found in log file"
}

func (m *logMatcher) NegatedFailureMessage(actual interface{}) string {
	return "Expected no log data to be found in log file"
}
