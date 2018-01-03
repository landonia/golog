package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/landonia/golog"
	"github.com/landonia/golog/prettylog"
	"github.com/landonia/golog/zerolog"
)

func main() {

	// Extract the log level
	logLevel := flag.String("loglevel", "debug", "Define the log level for the application")
	flag.Parse()

	// Set the global logging level
	golog.SetGlobalLevel(golog.LevelFromString(*logLevel))

	// Setup the log
	log, _ := prettylog.New(prettylog.WithNameSpace("golog.example"), prettylog.WithFlags(log.Ldate|log.Ltime|log.Lshortfile))

	// Test the basic logger
	logValues(log)

	// Setup a JSON based zero logger
	log, _ = zerolog.New(zerolog.WithNameSpace("golog.example"))

	// Now test the zerolog logger
	logValues(log)

	// Setup a pretty print based zero logger
	log, _ = zerolog.New(zerolog.WithPrettyConsole(), zerolog.WithNameSpace("golog.example"))

	// Now test the zerolog pretty logger
	logValues(log)
}

func logValues(log golog.Logger) {
	log.Info("Starting application....")
	log.Warn("Do not do that!")
	log.Debug("Sent %d value to server %s", 1, "google.com")
	log.Error("Error: %s", fmt.Errorf("Bang").Error())
	// log.Fatal("The application went boom")
}
