package loglevels

type LogLevel uint8

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
	PANIC
)

var logLevels = []string{
	"trace",
	"debug",
	"info",
	"warning",
	"error",
	"fatal",
	"panic",
}

func (logLevel LogLevel) String() string {
	return logLevels[logLevel]
}
