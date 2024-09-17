package logrus

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/SergeyIvanovDevelop/metrics-collection-and-alerting-service/pkg/logger/logiface"
	"github.com/SergeyIvanovDevelop/metrics-collection-and-alerting-service/pkg/logger/loglevels"
	"github.com/sirupsen/logrus"
)

type ExtLogrus struct {
	*logrus.Logger
}

type ExtEntry struct {
	*logrus.Entry
}

type FormatterType uint8

const (
	TextFormatterType FormatterType = iota
	JSONFormatterType
)

func New() *ExtLogrus {
	logrusLogger := logrus.New()
	return &ExtLogrus{
		Logger: logrusLogger,
	}
}

func Initialize(logOutput io.Writer, logLevel string, formatter FormatterType) *ExtLogrus {
	extLogrusLogger := New()
	extLogrusLogger.SetOutput(logOutput)
	extLogrusLogger.SetReportCaller(false)

	callerPrettyfier := func(f *runtime.Frame) (string, string) {
		filename := f.File
		line := f.Line
		funcName := f.Function
		shortFile := filename[strings.LastIndex(filename, "/")+1:]
		return funcName, shortFile + ":" + strconv.Itoa(line)
	}

	switch formatter {
	case TextFormatterType:
		extLogrusLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:    true,
			CallerPrettyfier: callerPrettyfier,
		})
	case JSONFormatterType:
		extLogrusLogger.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: callerPrettyfier,
			PrettyPrint:      true,
		})
	default:
		log.Fatalf("Неизвестный тип форматтера: %T\n", formatter)
	}

	switch logLevel {
	case loglevels.TRACE.String():
		extLogrusLogger.Level = logrus.TraceLevel
	case loglevels.DEBUG.String():
		extLogrusLogger.Level = logrus.DebugLevel
	case loglevels.INFO.String():
		extLogrusLogger.Level = logrus.InfoLevel
	case loglevels.WARNING.String():
		extLogrusLogger.Level = logrus.WarnLevel
	case loglevels.ERROR.String():
		extLogrusLogger.Level = logrus.ErrorLevel
	case loglevels.FATAL.String():
		extLogrusLogger.Level = logrus.FatalLevel
	case loglevels.PANIC.String():
		extLogrusLogger.Level = logrus.PanicLevel
	default:
		var errorMessage = fmt.Sprintf("Unknown 'logLevel': %s\n", logLevel)
		extLogrusLogger.Fatal(errorMessage)
	}

	return extLogrusLogger
}

func (extLogrus *ExtLogrus) WithFieldsIface(fields logiface.Fields) logiface.Logger {
	logrusEntry := extLogrus.WithFields(logrus.Fields(fields))
	extEntry := &ExtEntry{
		Entry: logrusEntry,
	}
	return extEntry
}

// ExtEntry implements logiface.Logger
func (extEntry *ExtEntry) WithFieldsIface(fields logiface.Fields) logiface.Logger {
	return extEntry
}

// ExtEntry implements logiface.Logger
func (extEntry *ExtEntry) WrapError(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// ExtLogrus implements logiface.Logger
func (extLogrus *ExtLogrus) WrapError(format string, args ...interface{}) error {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("не удалось получить информацию о вызове: %s", format)
	}

	shortFileName := file[strings.LastIndex(file, "/")+1:]

	function := runtime.FuncForPC(pc)
	shortFunctionName := "unknown"
	if function != nil {
		functionName := function.Name()
		shortFunctionName = functionName[strings.LastIndex(functionName, "/")+1:]
	}

	message := fmt.Sprintf("ошибка в функции '%s' (%s:%d): %s", shortFunctionName, shortFileName, line, format)
	return fmt.Errorf(message, args...)
}
