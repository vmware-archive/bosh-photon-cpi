package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

const (
	infoStr = "INFO "
	errStr  = "ERROR "
)

// Simple logger interface for reporting logs to bosh
type Logger interface {
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	// Closes underlying file
	Close()
	// Returns full filename of log
	Filename() string
}

type tempFileLogger struct {
	logFile *os.File
}

func New(logName string) (Logger, error) {
	dir := path.Join(os.TempDir(), "bosh-esxcloud-cpi-logs")
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return nil, err
	}
	file, err := ioutil.TempFile(dir, logName)
	if err != nil {
		return nil, err
	}
	return tempFileLogger{file}, nil
}

func (l tempFileLogger) Info(v ...interface{}) {
	l.logFile.WriteString(timestamp() + infoStr + fmt.Sprint(v...) + "\n")
}

func (l tempFileLogger) Infof(format string, v ...interface{}) {
	l.logFile.WriteString(timestamp() + infoStr + fmt.Sprintf(format, v...) + "\n")
}

func (l tempFileLogger) Error(v ...interface{}) {
	l.logFile.WriteString(timestamp() + errStr + fmt.Sprint(v...) + "\n")
}

func (l tempFileLogger) Errorf(format string, v ...interface{}) {
	l.logFile.WriteString(timestamp() + errStr + fmt.Sprintf(format, v...) + "\n")
}

func (l tempFileLogger) Close() {
	l.logFile.Close()
}

func (l tempFileLogger) Filename() string {
	return l.logFile.Name()
}

func timestamp() string {
	// UTC time formatted as RFC3339, retains order when sorted as a string
	return time.Now().UTC().Format(time.RFC3339) + " " // separator for parsing
}
