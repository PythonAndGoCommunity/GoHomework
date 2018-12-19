This is a simplified pure-Go redis.
It requires Docker, golang image and Go to install and launch.

1. How to build it:
	Move this folder to your Go workspace, then open the terminal and write command "make build". This will make a docker image with tag "custom_redis".

2. How to run it:
	After building this project, open the terminal and write command "make run". This will will create a docker container with tag "my_redis". It's entrypoint is set to be server part. To launch client part, open another terminal and write command "docker exec -it my_redis client [ARGUMENTS]".

3. How to check it:
	Open the terminal, go to the project's directory, type "make check".

4. How to test it:
	Open the terminal, go to the project's directory, type "make test". Coverage of server source code and client source code will be located in the project's directory under names "coverage_server.out" and "coverage_client.out" accordingly. 

5. If you want to change server launch arguments, open the Dockerfile and change ENTRYPOINT.

6. Available commands are:
		 - SET KEY VALUE - sets value to a key. If this key already exists, it will update the value. In case of success the answer will be either "OK.", or "Replaced existing value - PREVIOUS_VALUE" if key already exists.
		 - GET KEY - returns value, assigned to this key, or "(nil)"(notice, that it's a string, not the real nil value, cause nil can't be converted to string, so please be careful, when typing something similar as value), if key doesn't exist. 
		 - DEL KEY [KEY[...]] - deletes keys from database and returns, what keys were deleted and what were ignored, 'cause they don't exist. 

	6.1. There is an additional option of inputting arguments to commands implemented. You can use quotes to input keys and values, containing spaces(quotes will be dropped and the string inside quotes is considered). But, because of this, using quotes in keys and values themselves is forbidden. Also, "_" and " " are considered the same symbol.
				 Example: "value 1", value_1, "value_1" are considered the same.

P.S. Honestly, i had some troubles with importing in Go, because it supports only absolute path importing through usage of GOPATH and GOROOT variables, so imports in my project may look a little strange.
