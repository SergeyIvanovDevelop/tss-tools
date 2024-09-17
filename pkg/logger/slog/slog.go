package slog

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/SergeyIvanovDevelop/tss-tools/pkg/logger/logiface"
	"github.com/SergeyIvanovDevelop/tss-tools/pkg/logger/loglevels"
)

type ExtSlog struct {
	*slog.Logger
}

type FormatterType uint8

const (
	TextFormatterType FormatterType = iota
	JSONFormatterType
)

func New(handler slog.Handler) *ExtSlog {
	slogLogger := slog.New(handler)
	return &ExtSlog{
		Logger: slogLogger,
	}
}

func convertSlogToExtSlog(slogLogger *slog.Logger) *ExtSlog {
	return &ExtSlog{
		Logger: slogLogger,
	}
}

func Initialize(logOutput io.Writer, logLevel string, formatter FormatterType) *ExtSlog {
	var opts = &slog.HandlerOptions{}

	switch logLevel {
	case loglevels.TRACE.String():
		opts.Level = slog.LevelDebug
	case loglevels.DEBUG.String():
		opts.Level = slog.LevelDebug
	case loglevels.INFO.String():
		opts.Level = slog.LevelInfo
	case loglevels.WARNING.String():
		opts.Level = slog.LevelWarn
	case loglevels.ERROR.String():
		opts.Level = slog.LevelError
	case loglevels.FATAL.String():
		opts.Level = slog.LevelError
	case loglevels.PANIC.String():
		opts.Level = slog.LevelError
	default:
		var errorMessage = fmt.Sprintf("Unknown 'logLevel': %s\n", logLevel)
		log.Fatal(errorMessage)
	}

	opts.AddSource = false

	var h slog.Handler
	switch formatter {
	case TextFormatterType:
		h = slog.NewTextHandler(logOutput, opts)
	case JSONFormatterType:
		h = slog.NewJSONHandler(logOutput, opts)
	default:
		log.Fatalf("Неизвестный тип форматтера: %T\n", formatter)
	}

	extSlogLogger := New(h)

	return extSlogLogger
}

func (extSlog *ExtSlog) WithFieldsIface(fields logiface.Fields) logiface.Logger {
	var keyValues []any

	for k, v := range fields {
		keyValues = append(keyValues, k)
		keyValues = append(keyValues, v)
	}

	extSlogWith := convertSlogToExtSlog(extSlog.With(keyValues...))
	return extSlogWith
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) Info(args ...interface{}) {
	str := fmt.Sprintln(args...)
	extSlog.Logger.Info(str)
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) Infof(format string, args ...interface{}) {
	formattedString := fmt.Sprintf(format, args...)
	extSlog.Logger.Info(formattedString)
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) Warn(args ...interface{}) {
	str := fmt.Sprintln(args...)
	extSlog.Logger.Warn(str)
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) Warnf(format string, args ...interface{}) {
	formattedString := fmt.Sprintf(format, args...)
	extSlog.Logger.Warn(formattedString)
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) Debug(args ...interface{}) {
	str := fmt.Sprintln(args...)
	extSlog.Logger.Debug(str)
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) Debugf(format string, args ...interface{}) {
	formattedString := fmt.Sprintf(format, args...)
	extSlog.Logger.Debug(formattedString)
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) Trace(args ...interface{}) {
	str := fmt.Sprintln(args...)
	extSlog.Logger.Debug(str)
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) Tracef(format string, args ...interface{}) {
	formattedString := fmt.Sprintf(format, args...)
	extSlog.Logger.Debug(formattedString)
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) Error(args ...interface{}) {
	str := fmt.Sprintln(args...)
	extSlog.Logger.Error(str)
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) Errorf(format string, args ...interface{}) {
	formattedString := fmt.Sprintf(format, args...)
	extSlog.Logger.Error(formattedString)
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) Panic(args ...interface{}) {
	str := fmt.Sprintln(args...)
	extSlog.Logger.Error(str)
	panic(str)
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) Panicf(format string, args ...interface{}) {
	formattedString := fmt.Sprintf(format, args...)
	extSlog.Logger.Error(formattedString)
	panic(formattedString)
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) Fatal(args ...interface{}) {
	str := fmt.Sprintln(args...)
	extSlog.Logger.Error(str)
	os.Exit(1)
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) Fatalf(format string, args ...interface{}) {
	formattedString := fmt.Sprintf(format, args...)
	extSlog.Logger.Error(formattedString)
	os.Exit(1)
}

// ExtSlog implements logiface.Logger
func (extSlog *ExtSlog) WrapError(format string, args ...interface{}) error {
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
