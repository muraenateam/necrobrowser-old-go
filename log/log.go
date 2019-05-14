package log

import (
	ll "github.com/evilsocket/islazy/log"
)

// SetLevel defines the level of logging verbosity
func SetLevel(verbosity ll.Verbosity) {
	ll.Level = verbosity
}

func Debug(format string, args ...interface{}) {
	ll.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	ll.Info(format, args...)
}

func Important(format string, args ...interface{}) {
	ll.Important(format, args...)
}

func Warning(format string, args ...interface{}) {
	ll.Warning(format, args...)
}

func Error(format string, args ...interface{}) {
	ll.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	ll.Fatal(format, args...)
}
