package internal

import "log"

// Logger helps log output for benchdraw
type Logger struct {
	Verbosity int
	Logger    *log.Logger
}

// Log a message if verbosity <= Logger.Verbosity.  Uses Printf format.
func (l *Logger) Log(verbosity int, msg string, fmtArgs ...interface{}) {
	if l.Verbosity >= verbosity {
		l.Logger.Printf(msg, fmtArgs...)
	}
}
