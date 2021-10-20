package logger

type Handler interface {

	// Info logs a message at level Info.
	Info(args ...interface{})

	// Infof formats according to a format specifier and logs a message at level Info.
	Infof(format string, args ...interface{})

	// Warn logs a message at level Warn.
	Warn(args ...interface{})

	// Warnf formats according to a format specifier and logs a message at level Warn.
	Warnf(format string, args ...interface{})

	// Error logs a message at level Error.
	Error(args ...interface{})

	// Errorf formats according to a format specifier and logs a message at level Error.
	Errorf(format string, args ...interface{})

	// Debug logs a message at level Debug.
	Debug(args ...interface{})

	// Debugf formats according to a format specifier and logs a message at level Debug.
	Debugf(format string, args ...interface{})

	// Fatal logs a message at level Fatal
	// Then the process will exit with status 1
	Fatal(args ...interface{})

	// Fatalf formats according to a format specifier and logs a message at level Fatal.
	// Then the process will exit with status 1
	Fatalf(format string, args ...interface{})
}
