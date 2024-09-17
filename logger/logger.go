package logger

import (
	"fmt"
	"io"

	"github.com/SergeyIvanovDevelop/metrics-collection-and-alerting-service/pkg/logger/logiface"
	// logrusext "github.com/SergeyIvanovDevelop/metrics-collection-and-alerting-service/pkg/logger/logrus"
	slogext "github.com/SergeyIvanovDevelop/metrics-collection-and-alerting-service/pkg/logger/slog"
)

type Fields logiface.Fields
type Log logiface.Logger

// Log is 'singleton' pattern
var Logger logiface.Logger

func Initialize(logOutput io.Writer, logLevel string) {
	// Logger = logrusext.Initialize(logOutput, logLevel, logrusext.TextFormatterType)
	Logger = slogext.Initialize(logOutput, logLevel, slogext.TextFormatterType)
}

func Trace(args ...interface{}) {
	Logger.Trace(args)
}

func Tracef(format string, args ...interface{}) {
	Logger.Tracef(format, args)
}

func Debug(args ...interface{}) {
	Logger.Debug(args)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args)
}

func Info(args ...interface{}) {
	Logger.Info(args)
}

func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args)
}

func Warn(args ...interface{}) {
	Logger.Warn(args)
}

func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args)
}

func Error(args ...interface{}) {
	Logger.Error(args)
}

func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args)
}

func Fatal(args ...interface{}) {
	Logger.Error(args)
}

func Fatalf(format string, args ...interface{}) {
	Logger.Errorf(format, args)
}

func WithFields(fields Fields) logiface.Logger {
	return Logger.WithFieldsIface(logiface.Fields(fields))
}

func WrapError(format string, args ...interface{}) error {
	return Logger.WrapError(format, args)
}

// For packages loggers
func lazyPackageLoggerInit(pkgLog logiface.Logger, pkgName string) logiface.Logger {
	if pkgLog == nil {
		if Logger != nil {
			pkgLog = WithFields(Fields{
				"package": pkgName,
			})
		} else {
			panicString := fmt.Sprintf("[logger] lazyPackageLoggerInit for package '%s' failed", pkgName)
			panic(panicString)
		}
	}
	return pkgLog
}

func AddLoggerFields(pkgLog logiface.Logger, pkgName string, fields Fields) logiface.Logger {
	pkgLog = lazyPackageLoggerInit(pkgLog, pkgName)
	return pkgLog.WithFieldsIface(logiface.Fields(fields))
}
