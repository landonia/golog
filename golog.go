// Copyright 2018 Landon Wainwright. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golog

import (
	"strings"
	"sync/atomic"
)

// LevelNotifierHandler is used to notify any interested golog.Logger impls
// that the global level has changed
type LevelNotifierHandler func(Level)

var (
	// globalLogLevel is the current global logging level
	globalLogLevel = new(uint32)

	// Hold the handlers
	handlers []LevelNotifierHandler
)

// LevelChangeNotifier allows a handler to be added to be notified of any global
// level change
func LevelChangeNotifier(handler LevelNotifierHandler) {
	handlers = append(handlers, handler)
}

// SetGlobalLevel will set the global log level to the specified level
// If the log level is not recogised it will return a false and default to INFO
func SetGlobalLevel(level Level) {
	atomic.StoreUint32(globalLogLevel, uint32(level))
	for _, handler := range handlers {
		handler(level)
	}
}

func GlobalLevel() Level {
	return Level(atomic.LoadUint32(globalLogLevel))
}

// Level will define a log level
type Level uint8

// Define the logging levels
const (

	// No level has been defined
	NONE Level = iota

	// DISABLED will disable the logger logging level
	DISABLED

	// FATAL logging level
	FATAL

	// ERROR logging level
	ERROR

	// WARN logging level
	WARN

	// INFO logging level
	INFO

	// DEBUG logging level
	DEBUG

	// TRACE logging level
	TRACE
)

func (l Level) String() string {
	switch l {
	case TRACE:
		return "TRACE"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return ""
}

// LevelFromString will convert the string to a Level
// If the string cannot be converted, it will return NONE
func LevelFromString(ll string) Level {
	switch strings.ToUpper(ll) {
	case "FATAL":
		return FATAL
	case "ERROR":
		return ERROR
	case "WARN":
		return WARN
	case "INFO":
		return INFO
	case "DEBUG":
		return DEBUG
	case "TRACE":
		return TRACE
	case "DISABLED":
		return DISABLED
	}
	return NONE
}

// Logger defines a basic set of logging functions
// For formatting refer to the underlying log implementation
type Logger interface {

	// SetLevel will set the logging level for this logger
	SetLevel(level Level)

	// Fatal prints a Fatal level message
	Fatal(format string, params ...interface{})

	// Error prints an Error level message
	Error(format string, params ...interface{})

	// Warn prints a Warn level message
	Warn(format string, params ...interface{})

	// Info prints an Info level message
	Info(format string, params ...interface{})

	// Debug prints a Debug level message
	Debug(format string, params ...interface{})

	// Trace prints a Trace level message
	Trace(format string, params ...interface{})
}

// EmptyLog can be used to fill the interface requirement where a logger
// has not been provided
type EmptyLogger struct{}

func (el *EmptyLogger) SetLevel(level Level)                       {}
func (el *EmptyLogger) Fatal(format string, params ...interface{}) {}
func (el *EmptyLogger) Error(format string, params ...interface{}) {}
func (el *EmptyLogger) Warn(format string, params ...interface{})  {}
func (el *EmptyLogger) Info(format string, params ...interface{})  {}
func (el *EmptyLogger) Debug(format string, params ...interface{}) {}
func (el *EmptyLogger) Trace(format string, params ...interface{}) {}
