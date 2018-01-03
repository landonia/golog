// Copyright 2018 Landon Wainwright. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prettylog

// prettylog.go provides a very simple console based logger also adding in
// some nice pretty console strings. It is not massively efficient as it
// is just using

import (
	"log"
	"os"

	"github.com/landonia/golog"
	"github.com/landonia/util"
	"github.com/landonia/util/color"
)

// WithConfiguration defines a function type to configure the the logger
type WithConfiguration func(*PrettyLogger) error

// WithNameSpace will return a WithConfiguration function to define a namespace
// within the logger
func WithNameSpace(ns string) WithConfiguration {
	return func(logger *PrettyLogger) error {
		logger.ns = ns
		return nil
	}
}

// WithFlags will change the underlying log.Logger.flags field which defines
// what is pre-pended to the output (see https://golang.org/pkg/log/#pkg-constants)
func WithFlags(flag int) WithConfiguration {
	return func(logger *PrettyLogger) error {
		logger.out.SetFlags(flag)
		return nil
	}
}

// WithLevel will set the logger level
func WithLevel(level golog.Level) WithConfiguration {
	return func(logger *PrettyLogger) error {
		logger.SetLevel(level)
		return nil
	}
}

// WithFileOutput will override the log from printing to stdout and instead print to the specified file
// An error will be returned if the file could not be opened or created
func WithFileOutput(filename string) WithConfiguration {
	return func(logger *PrettyLogger) error {

		// Attempt to open the file (create if it does not exist) and put the file in append mode
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {

			// The log could not be found or created
			return err
		}

		// Set the output to the file
		logger.file = f
		logger.out.SetOutput(f)
		return nil
	}
}

// // Close will close the underlying file
// func Close() error {
// 	if file != nil {
// 		return file.Close()
// 	}
// 	return nil
// }

// PrettyLogger is a simple logger providing a namespace
// and providing granular logging levels
type PrettyLogger struct {
	file  *os.File    // file to write the log to
	out   *log.Logger // out is the base logger and can be overwritten on a package level if required
	level golog.Level // The level for this logger
	ns    string      // The namespace for this log
}

// New will return a new PrettyLogger for the particular namespace
func New(withConfs ...WithConfiguration) (golog.Logger, error) {
	pl := &PrettyLogger{}

	// Setup a default logger
	pl.out = log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Add the WithConfiguration functions
	for _, withConf := range withConfs {
		if err := withConf(pl); err != nil {
			return nil, err
		}
	}
	return pl, nil
}

// print message to standard out prefixed with date and time
func (gl *PrettyLogger) print(level golog.Level, ns, s string) {
	if level <= gl.level || level <= golog.GlobalLevel() {
		l := level.String()
		switch level {
		case golog.FATAL:
		case golog.ERROR:
			l = color.PrintRedMaxWidth(level, 5)
		case golog.WARN:
			l = color.PrintYellowMaxWidth(level, 5)
		case golog.INFO:
			l = color.PrintGreenMaxWidth(level, 5)
		}
		gl.out.Print(util.FormatString("[%-5s] [%s] %s", l, ns, s))
	}
}

// SetLevel will define the level of logging for just this log (and sublogs,
// unless overwritten)
func (gl *PrettyLogger) SetLevel(level golog.Level) {
	gl.level = level
}

// PrintColour prints colored message
func (gl *PrettyLogger) PrintColour(level golog.Level, s string, c color.Color) {
	gl.print(level, gl.ns, color.PrintColor(s, c))
}

// Fatal prints a Fatal level message
func (gl *PrettyLogger) Fatal(format string, params ...interface{}) {
	gl.PrintColour(golog.FATAL, util.FormatString(format, params...), color.RED)
	os.Exit(1)
}

// Error prints an Error level message
func (gl *PrettyLogger) Error(format string, params ...interface{}) {
	gl.PrintColour(golog.ERROR, util.FormatString(format, params...), color.RED)
}

// Warn prints a Warn level message
func (gl *PrettyLogger) Warn(format string, params ...interface{}) {
	gl.PrintColour(golog.WARN, util.FormatString(format, params...), color.YELLOW)
}

// Info prints an Info level message
func (gl *PrettyLogger) Info(format string, params ...interface{}) {
	gl.PrintColour(golog.INFO, util.FormatString(format, params...), color.GREEN)
}

// Debug prints a Debug level message
func (gl *PrettyLogger) Debug(format string, params ...interface{}) {
	gl.print(golog.DEBUG, gl.ns, util.FormatString(format, params...))
}

// Trace prints a Trace level message
func (gl *PrettyLogger) Trace(format string, params ...interface{}) {
	gl.print(golog.TRACE, gl.ns, util.FormatString(format, params...))
}
