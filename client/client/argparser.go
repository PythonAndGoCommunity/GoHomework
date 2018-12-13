package client

import "flag"

func GetCommandLineParams() (string, string, bool, string) {
	var host, port, filename string
	var dump bool

	flag.StringVar(&host, "host", "localhost", "default connection host")
	flag.StringVar(&host, "h", "localhost", "default connection host (shortcut)")
	flag.StringVar(&port, "port", "9090", "connection port")
	flag.StringVar(&port, "p", "9090", "connection port (shortcut)")
	flag.BoolVar(&dump, "dump", false, "dump database into JSON format")
	flag.BoolVar(&dump, "d", false, "dump database into JSON format (shortcut)")
	flag.StringVar(&filename, "restore", "", "restore the database from the dumped file")
	flag.StringVar(&filename, "r", "", "restore the database from the dumper file (shortcut)")
	flag.Parse()

	return host, port, dump, filename
}
