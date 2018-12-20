# GoHomework
Implementation of the server-client solution for storing KV data lightweight analog of the Redis (https://redis.io/).
More in file Task.pdf

use makefile for all actions

make check:
    run  "go vet", "goimports", "golint" 

make test:
    run tests and save the output to the "coverage.out" file.

Work with application

Server
You may run server using arguments:
    -p --port 				listening port, default is 9090;
    -m --mode 				storage mode, default is "memory", alternate mode is "disk" (save to file "data.json")
	-v --verbose 			verbose mode, full log of the client requests

Client
You may run client using arguments:
    -p --port 				connect to port, default is 9090;
    -h --host 				connect to ip address, default is 127.0.0.1
	--dump					dump the whole database to the JSON format on STDOUT (example:'[{"key": "YOUR_KEY", "value": "YOUR_VALUE"}]'). Save to file 'data.json'
	--restore				restore the database from the dumped file 'data.json'.

Commands:
    set key value			updates one key at a time with the given value
    get key                 returns tuple of the value and the key state. The state either present or absent
    del key                 removes one key at a time and returns the state of the resource
    keys [searching key]    returns all keys matching pattern, for example "h?llo" matches "hello", "hallo" and "hxllo"
	exit                    exit