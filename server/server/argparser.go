package server

import (
	"flag"
)

func GetCommandLineParams() (string, string, bool) {
	var storage, port string
	var verbose bool

	flag.StringVar(&storage, "mode", "disk", "possible storage option")
	flag.StringVar(&storage, "m", "disk", "possible storage option (shortcut)")
	flag.StringVar(&port, "port", "9090", "server default port")
	flag.StringVar(&port, "p", "9090", "server default port (shortcut)")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode, full log of the client requests")
	flag.BoolVar(&verbose, "v", false, "verbose mode, full log of the client requests (shortcut)")
	flag.Parse()

	return storage, port, verbose
}
