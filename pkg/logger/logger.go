package logger

type Logger interface {
	LogLine(logLine string)
	StartLogging()
	StopLogging()
	LoggingEnabled() bool
}
