package logs

var (
	_logger Logger
)

// Debug debug log
func Debug(args ...interface{}) {
	_logger.Debug(args)
}

// Debugf info format log
func Debugf(format string, args ...interface{}) {
	_logger.Debugf(format, args)
}

// Info info log
func Info(args ...interface{}) {
	_logger.Info(args)
}

// Infof info format log
func Infof(format string, args ...interface{}) {
	_logger.Infof(format, args)
}

// Warn warn log
func Warn(args ...interface{}) {
	_logger.Warn(args)
}

// Warnf warn format log
func Warnf(format string, args ...interface{}) {
	_logger.Warnf(format, args)
}

// Error error log
func Error(args ...interface{}) {
	_logger.Error(args)
}

// Errorf error format log
func Errorf(format string, args ...interface{}) {
	_logger.Errorf(format, args)
}

// Panic panic log
func Panic(args ...interface{}) {
	_logger.Panic(args)
}

// Panicf panic format log
func Panicf(format string, args ...interface{}) {
	_logger.Panicf(format, args)
}

// Fatal log message, and call os.Exist.
func Fatal(args ...interface{}) {
	_logger.Fatal(args)
}

// Fatalf log format message, and call os.Exist.
func Fatalf(format string, args ...interface{}) {
	_logger.Fatalf(format, args)
}

// New create a logger.
func New(fn func() Logger) {
	if fn == nil {
		panic("[logs]: register logger is nil.")
	}

	_logger = fn()
}
