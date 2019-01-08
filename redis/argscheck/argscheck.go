package argscheck

//Start - update gPort, gMode, gIP variable if it's need to be changed
func Start(args []string, gPort string, gMemory string, gIP string) (string, string, string) {
	var (
		port = false
		host = false
		mode = false
	)
	for _, value := range args[1:] {
		if port {
			gPort = ":" + value
			port = false
		} else {
			if host {
				gIP = value
				host = false
			}
		}
		if mode {
			gMemory = value
			mode = false
		}

		switch value {
		case "-p", "--port":
			port = true
		case "-h", "--host":
			host = true
		case "-m", "--mode":
			mode = true
		}
	}
	return gPort, gIP, gMemory
}
