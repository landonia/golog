// Copyright 2016 Landon Wainwright. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golog

import (
	"fmt"
	"log"
	"os"
)

// Define the logging levels
const (

	// OFF logging level
	OFF int = 0

	// FATAL logging level
	FATAL int = 1

	// ERROR logging level
	ERROR int = 2

	// WARN logging level
	WARN int = 3

	// INFO logging level
	INFO int = 4

	// DEBUG logging level
	DEBUG int = 5

	// TRACE logging level
	TRACE int = 6
)

// Colour allows colours to be defined
type Colour int

// Define some colours for the output text
const (

	// RED color
	RED Colour = 31

	// GREEN color
	GREEN Colour = 32

	// YELLOW color
	YELLOW Colour = 33

	// BLUE color
	BLUE Colour = 34

	// MAG color
	MAG Colour = 35

	// CYAN color
	CYAN Colour = 36
)

var (

	// LoggingLevel is the current global logging level
	LoggingLevel = INFO

	// OutputLog is the base logger and can be overwritten on a package level if required
	OutputLog = log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// file to write the log to
	file *os.File
)

// OutputToFile will override the log from printing to stdout and instead print to the specified file
// An error will be returned if the file could not be opened or created
func OutputToFile(filename string) error {

	// Attempt to open the file (create if it does not exist) and put the file in append mode
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {

		// The log could not be found or created
		return err
	}

	// Set the output to the file
	file = f
	OutputLog.SetOutput(f)
	return nil
}

// GoLog is a wrapper
type GoLog struct {
	ns string // The namespace for this log
}

// New will return a new log for the particular namespace
func New(ns string) *GoLog {
	return &GoLog{ns}
}

// Close will close the underlying file
func Close() error {
	if file != nil {
		return file.Close()
	}
	return nil
}

// print message to standard out prefixed with date and time
func print(level int, ns, s string) {
	if LoggingLevel >= level {
		l := "TRACE"
		switch level {
		case FATAL:
			l = "FATAL"
		case ERROR:
			l = "ERROR"
		case WARN:
			l = "WARN"
		case INFO:
			l = "INFO"
		case DEBUG:
			l = "DEBUG"
		}
		OutputLog.Print(fmt.Sprintf("[%s] [%s] %s", l, ns, s))
	}
}

// PrintColour prints coloured message
func (gl *GoLog) PrintColour(level int, s string, colour Colour) {
	print(level, gl.ns, fmt.Sprintf("\x1b[%v;1m%v\x1b[0m", colour, s))
}

// Fatal prints a Fatal level message
func (gl *GoLog) Fatal(s string) {
	gl.PrintColour(FATAL, s, RED)
}

// Error prints an Error level message
func (gl *GoLog) Error(s string) {
	gl.PrintColour(ERROR, s, RED)
}

// Warn prints a Warn level message
func (gl *GoLog) Warn(s string) {
	gl.PrintColour(WARN, s, YELLOW)
}

// Info prints an Info level message
func (gl *GoLog) Info(s string) {
	gl.PrintColour(INFO, s, GREEN)
}

// Debug prints a Debug level message
func (gl *GoLog) Debug(s string) {
	print(DEBUG, gl.ns, s)
}

// Trace prints a Trace level message
func (gl *GoLog) Trace(s string) {
	print(TRACE, gl.ns, s)
}
