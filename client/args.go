package main

import (
	"flag"
)

// Config represents a set of parsed flags for client settings.
type Config struct {
	port    int
	host    string
	dump    bool
	restore bool
}

var config Config

func init() {
	const (
		defaultPort = 9090
		portDesc    = "The port to connect to the server"

		defaultHost = "127.0.0.1"
		hostDesc    = "The host to connect to the server"

		defaultDump = false
		dumpDesc    = "Dump the whole database to the JSON format on STDOUT"

		defaultRestore = false
		restoreDesc    = "Restore the database from the dumped file"
	)

	shorhand := " (shorthand)"

	flag.IntVar(&config.port, "port", defaultPort, portDesc)
	flag.IntVar(&config.port, "p", defaultPort, portDesc+shorhand)

	flag.StringVar(&config.host, "host", defaultHost, hostDesc)
	flag.StringVar(&config.host, "h", defaultHost, hostDesc+shorhand)
}
