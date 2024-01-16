package log

import (
	"github.com/ralstan-vaz/go-errors"
	"stash.bms.bz/merchandise/go-logger"
)

const (
	// Priority1 ...
	Priority1 string = "priority 1"
	// Priority2 ...
	Priority2 string = "priority 2"
)

// Logger is a global singleton
var Logger *logger.Logger

// InitLogger ...
func InitLogger() {
	// Initializes logger
	// inits with app name
	// caller level is 2 so that we get info from where the outer function was called
	newLogger := logger.New("go-boilerplate", 2)
	// makes logger a global singleton
	logger.Global(newLogger)
	// enables debug logs
	logger.GLog.Set(`{"debug":true,"reference":"string"}`)

	Logger = logger.GLog
}

// Info ...
func Info(description string, reference ...interface{}) {
	if Logger == nil {
		return
	}

	Logger.Info(description, reference...)
}

// Debug ...
func Debug(description string, reference ...interface{}) {
	if Logger == nil {
		return
	}

	Logger.Debug(description, reference...)
}

// Warn ...
func Warn(description string, reference ...interface{}) {
	if Logger == nil {
		return
	}

	Logger.Warn(description, reference...)
}

// Error ...
func Error(code string, description string, severity string, source interface{}, reference ...interface{}) {
	if Logger == nil {
		return
	}

	Logger.Error(code, description, severity, mapSource(source), reference...)
}

// Fatal ...
func Fatal(code string, description string, source interface{}, reference ...interface{}) {
	if Logger == nil {
		return
	}

	Logger.Fatal(code, description, mapSource(source), reference...)
}

// used to map a different error source to the logger error source struct
func mapSource(source interface{}) *logger.ErrorSource {
	var errSource *logger.ErrorSource
	errSource = new(logger.ErrorSource)

	newSource, ok := source.(errors.Source)
	if !ok {
		errSource = nil
	} else {
		errSource.Caller = &newSource.Caller
		errSource.File = &newSource.File
		errSource.Line = &newSource.Line
		errSource.StackTrace = &newSource.StackTrace
		if newSource.Error != nil {
			errSource.ErrorSource = newSource.Error
			errStr := newSource.Error.Error()
			errSource.ErrorString = &errStr
		}
	}

	return errSource
}
