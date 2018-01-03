// Copyright 2016 Landon Wainwright. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zerolog

// zerolog.go wraps the zerolog.Logger to adhere to the golog.Logger interface

import (
	"os"

	"github.com/landonia/golog"
	"github.com/landonia/util"
	"github.com/rs/zerolog"
)

// Register for any global stats changed
func init() {
	// Register a handler for the global level change
	golog.LevelChangeNotifier(setGlobalLevel)
}

// WithConfiguration defines a function type to configure the the logger
type WithConfiguration func(*ZeroLogger) error

// WithNameSpace will return a WithConfiguration function to define a namespace
// within the logger
func WithNameSpace(ns string) WithConfiguration {
	return func(logger *ZeroLogger) error {
		logger.zl = logger.zl.With().Str("ns", ns).Logger()
		return nil
	}
}

// WithPrettyConsole will return a WithConfiguration function that will use
// the pretty logger making zerologger print pretty strings to the output and
// not use JSON
func WithPrettyConsole() WithConfiguration {
	return func(logger *ZeroLogger) error {
		// Update the output to use the pretty logger
		logger.zl = logger.zl.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		return nil
	}
}

// WithLevel will set the logger level
func WithLevel(level golog.Level) WithConfiguration {
	return func(logger *ZeroLogger) error {
		logger.SetLevel(level)
		return nil
	}
}

// WithFileOutput will override the log from printing to stdout and instead print to the specified file
// An error will be returned if the file could not be opened or created
func WithFileOutput(filename string) WithConfiguration {
	return func(logger *ZeroLogger) error {

		// Attempt to open the file (create if it does not exist) and put the file in append mode
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {

			// The log could not be found or created
			return err
		}

		// Set the output to the file
		logger.file = f
		logger.zl = logger.zl.Output(f)
		return nil
	}
}

// ZeroLogger is a wrapper for the zerolog.Logger type using
// the golog.Logger interface
type ZeroLogger struct {
	zl   zerolog.Logger // The logger
	file *os.File       // file to write the log to
}

// setGlobalLevel will take the golog level and match it to the equivalent
// zerolog level
func setGlobalLevel(level golog.Level) {
	setLoggerLevel(func(lvl zerolog.Level) {
		zerolog.SetGlobalLevel(lvl)
	}, level)
}

// setLoggerLevel will map the golog.Level to the root logger level type
func setLoggerLevel(set func(zerolog.Level), level golog.Level) {
	switch level {
	case golog.FATAL:
		set(zerolog.FatalLevel)
	case golog.ERROR:
		set(zerolog.ErrorLevel)
	case golog.WARN:
		set(zerolog.WarnLevel)
	case golog.DEBUG:
	case golog.TRACE:
		set(zerolog.DebugLevel)
	case golog.DISABLED:
		set(zerolog.Disabled)
	default:
		set(zerolog.InfoLevel)
	}
}

// NewZeroLogger will return a new log for the particular namespace
// wrapping the zerolog implementation (for speed)
func New(withConfs ...WithConfiguration) (golog.Logger, error) {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return newZeroLogger(log, withConfs...)
}

// SetLevel will define the level of logging for just this log (and sublogs,
// unless overwritten)
func (log *ZeroLogger) SetLevel(level golog.Level) {
	setLoggerLevel(func(lvl zerolog.Level) {
		log.zl.Level(lvl)
	}, level)
}

// SubLogger will create a new sublogger and creates a new context that is used
// for this logger (plus any sub loggers)
func (log *ZeroLogger) SubLogger(withConfs ...WithConfiguration) (golog.Logger, error) {
	return newZeroLogger(log.zl.With().Logger(), withConfs...)
}

// newZeroLogger will return a new ZeroLogger instance pointer
func newZeroLogger(zl zerolog.Logger, withConfs ...WithConfiguration) (golog.Logger, error) {
	xl := &ZeroLogger{
		zl: zl,
	}
	// Add the WithConfiguration functions
	for _, withConf := range withConfs {
		if err := withConf(xl); err != nil {
			return nil, err
		}
	}
	return xl, nil
}

// Fatal prints a Fatal level message
func (log *ZeroLogger) Fatal(format string, params ...interface{}) {
	log.zl.Fatal().Msg(util.FormatString(format, params...))
}

// Error prints an Error level message
func (log *ZeroLogger) Error(format string, params ...interface{}) {
	log.zl.Error().Msg(util.FormatString(format, params...))
}

// Warn prints a Warn level message
func (log *ZeroLogger) Warn(format string, params ...interface{}) {
	log.zl.Warn().Msg(util.FormatString(format, params...))
}

// Info prints an Info level message
func (log *ZeroLogger) Info(format string, params ...interface{}) {
	log.zl.Info().Msg(util.FormatString(format, params...))
}

// Debug prints a Debug level message
func (log *ZeroLogger) Debug(format string, params ...interface{}) {
	log.zl.Debug().Msg(util.FormatString(format, params...))
}

// Trace prints a Trace level message
func (log *ZeroLogger) Trace(format string, params ...interface{}) {

	// There is no Trace level debugging in zerolog so we will revert to
	// the next level (which is Debug)
	log.Debug(util.FormatString(format, params...))
}
