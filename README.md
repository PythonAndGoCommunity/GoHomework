# Redis like database
Small Redis like database in Go.
Supported commands: SET, GET, DEL, KEYS
# Installation Instructions
Type `make` to build everything necessary
Use `src/server/server` for server and `src/client/client` for client.
Server must be run first.
# Running Redis like database in Docker
Run `docker build -t "gohomework" .` to build docker image.
Run `docker run gohomework` to run the server.

