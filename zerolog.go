// Copyright 2016 Landon Wainwright. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golog

import (
	"os"

	"github.com/rs/zerolog"
)

// GoZeroLog is a wrapper for the zerolog impl
type GoZeroLog struct {
	ns      string // The namespace for this log
	zl      zerolog.Logger
	Context zerolog.Context
}

// SetZeroGlobalLevel will take the golog level and match it to the equivalent
// zerolog level
func SetZeroGlobalLevel(level int) {
	setLoggerLevel(func(lvl zerolog.Level) {
		zerolog.SetGlobalLevel(lvl)
	}, level)
}

func setLoggerLevel(set func(zerolog.Level), level int) {
	switch level {
	case FATAL:
		set(zerolog.FatalLevel)
	case ERROR:
		set(zerolog.ErrorLevel)
	case WARN:
		set(zerolog.WarnLevel)
	case DEBUG:
	case TRACE:
		set(zerolog.DebugLevel)
	default:
		set(zerolog.InfoLevel)
	}
}

// NewZeroLogger will return a new log for the particular namespace
// wrapping the zerolog implementation (for speed)
func NewZeroLogger(ns string) *GoZeroLog {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return newZeroLogger(ns, log.With(), log)
}

// NewZeroLogger will return a new log for the particular namespace
// wrapping the zerolog implementation (for speed)
func NewPrettyZeroLogger(ns string) *GoZeroLog {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	return newZeroLogger(ns, log.With(), log)
}

// SetLoggerLevel will define the level of logging for just this log (and sublogs,
// unless overwritten)
func (log *GoZeroLog) SetLoggerLevel(level int) {
	setLoggerLevel(func(lvl zerolog.Level) {
		log.zl.Level(lvl)
	}, level)
}

// SubLogger will create a new sublogger and creates a new context that is used
// for this logger (plus any sub loggers)
func (log *GoZeroLog) SubLogger(ns string) *GoZeroLog {
	ctx := log.zl.With()
	return newZeroLogger(ns, ctx, ctx.Logger())
}

// newZeroLogger will return a new GoZeroLog instance pointer
func newZeroLogger(ns string, ctx zerolog.Context, zl zerolog.Logger) *GoZeroLog {
	return &GoZeroLog{
		ns:      ns,
		zl:      zl,
		Context: ctx,
	}
}

// Fatal prints a Fatal level message
func (log *GoZeroLog) Fatal(format string, params ...interface{}) {
	log.zl.Fatal().Msg(formatString(format, params...))
}

// Error prints an Error level message
func (log *GoZeroLog) Error(format string, params ...interface{}) {
	log.zl.Error().Msg(formatString(format, params...))
}

// Warn prints a Warn level message
func (log *GoZeroLog) Warn(format string, params ...interface{}) {
	log.zl.Warn().Msg(formatString(format, params...))
}

// Info prints an Info level message
func (log *GoZeroLog) Info(format string, params ...interface{}) {
	log.zl.Info().Msg(formatString(format, params...))
}

// Debug prints a Debug level message
func (log *GoZeroLog) Debug(format string, params ...interface{}) {
	log.zl.Debug().Msg(formatString(format, params...))
}

// Trace prints a Trace level message
func (log *GoZeroLog) Trace(format string, params ...interface{}) {

	// There is no Trace level debugging in zerolog so we will revert to
	// the next level (which is Debug)
	log.Debug(formatString(format, params...))
}
