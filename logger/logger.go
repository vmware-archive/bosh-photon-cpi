package logger

import (
	"bytes"
	"fmt"
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

	LogData() string
}

type bufferLogger struct {
	buffer *bytes.Buffer
}

func New() Logger {
	return bufferLogger{&bytes.Buffer{}}
}

func (l bufferLogger) Info(v ...interface{}) {
	l.buffer.WriteString(timestamp() + infoStr + fmt.Sprint(v...) + "\n")
}

func (l bufferLogger) Infof(format string, v ...interface{}) {
	l.buffer.WriteString(timestamp() + infoStr + fmt.Sprintf(format, v...) + "\n")
}

func (l bufferLogger) Error(v ...interface{}) {
	l.buffer.WriteString(timestamp() + errStr + fmt.Sprint(v...) + "\n")
}

func (l bufferLogger) Errorf(format string, v ...interface{}) {
	l.buffer.WriteString(timestamp() + errStr + fmt.Sprintf(format, v...) + "\n")
}

func (l bufferLogger) LogData() string {
	return l.buffer.String()
}

func timestamp() string {
	// UTC time formatted as RFC3339, retains order when sorted as a string
	return time.Now().UTC().Format(time.RFC3339) + " " // separator for parsing
}
