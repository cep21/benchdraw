package internal

import "log"

type Logger struct {
	Verbosity int
	Logger    *log.Logger
}

func (l *Logger) Log(verbosity int, msg string, fmtArgs ...interface{}) {
	if l.Verbosity >= verbosity {
		l.Logger.Printf(msg, fmtArgs...)
	}
}
