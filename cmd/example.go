package main

import (
	"flag"
	"fmt"

	"github.com/landonia/golog"
)

var (

	// Setup the log
	log = golog.New("golog.example")

	// Setup a JSON based zero logger
	zlog = golog.NewZeroLogger("golog.example")

	// Setup a pretty print based zero logger
	pzlog = golog.NewPrettyZeroLogger("golog.example")

	// Some example values
	host = "google.com"
)

func main() {

	// Extract the log level
	logLevel := flag.String("loglevel", "debug", "Define the log level for the application")
	flag.Parse()

	// Define the log level for the application
	if !golog.LogLevel(*logLevel) {
		log.Fatal("Incorrect log level: %s", *logLevel)
	}

	// Test the basic logger
	logValues(log)

	// Now test the zerolog logger
	logValues(zlog)

	// Now test the zerolog pretty logger`
	logValues(pzlog)
}

func logValues(log golog.GoLogger) {
	log.Info("Starting application....")
	log.Warn("Do not do that!")
	log.Debug("Sent %d value to server %s", 1, host)
	log.Error("Error: %s", fmt.Errorf("Bang").Error())
	// log.Fatal("The application went boom")
}
