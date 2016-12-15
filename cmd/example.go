package main

import (
	"flag"
	"fmt"

	"github.com/landonia/golog"
)

var (

	// Setup the log
	log = golog.New("golog.example")

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

	log.Info("Starting application....")
	log.Warn("Do not do that!")
	log.Debug("Sent %d value to server %s", 1, host)
	log.Error("Error: %s", fmt.Errorf("Bang").Error())
	log.Fatal("The application went boom")
}
