package main

import (
	"flag"
)

// Config represents a set of parsed flags.
type Config struct {
	port    int
	mode    string
	verbose bool
}

var config Config

func init() {
	const (
		defaultPort = 9090
		portDesc    = "The port for listening on"

		defaultMode = "disk"
		modeDesc    = "The possible storage option"

		defaultVerbose = false
		verboseDesc    = "Verbose mode, full log of the client requests"
	)

	shorhand := " (shorthand)"

	flag.IntVar(&config.port, "port", defaultPort, portDesc)
	flag.IntVar(&config.port, "p", defaultPort, portDesc+shorhand)

	flag.StringVar(&config.mode, "mode", defaultMode, modeDesc)
	flag.StringVar(&config.mode, "m", defaultMode, modeDesc+shorhand)

	flag.BoolVar(&config.verbose, "verbose", defaultVerbose, verboseDesc)
	flag.BoolVar(&config.verbose, "v", defaultVerbose, verboseDesc+shorhand)
}
