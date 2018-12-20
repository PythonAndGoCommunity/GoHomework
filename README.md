# GoHomework
Implementation of the server-client solution for storing KV data lightweight analog of the Redis (https://redis.io/).
More in file [Task.pdf](https://github.com/Dubouski/GoHomework/blob/working-branch/Task.pdf)

Use makefile for all actions

## Testing:
**make check**
Run  "go vet", "goimports", "golint" 

**make test**
Run tests and save the output to the "coverage.out" file.


## Work with application

## Server
You may run server using arguments:
-	'-p' or '--port' 				
>listening port, default is 9090;
- 	'-m' or '--mode' 				
>storage mode, default is "memory", alternate mode is "disk" (save to file "data.json");
- 	'-v' or '--verbose' 			
>verbose mode, full log of the client requests.

## Client
You may run client using arguments:
- 	'-p' or '--port'
>connect to port, default is 9090;
- 	'-h' or '--host'
>connect to ip address, default is 127.0.0.1;
- 	'--dump'
>dump the whole database to the JSON format on STDOUT (example:'[{"key": "YOUR_KEY", "value": "YOUR_VALUE"}]'). Save to file 'data.json';
- 	'--restore'
>restore the database from the dumped file 'data.json'.

## Commands:
updates one key at a time with the given value:
-     set key value
returns tuple of the value and the key state. The state either present or absent:
-     get key
removes one key at a time and returns the state of the resource:
-     del key 
returns all keys matching pattern, for example "h?llo" matches "hello", "hallo" and "hxllo":
-     keys [pattern]
exit from app
-     exit                   


