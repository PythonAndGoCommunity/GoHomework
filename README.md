Server-client solution for storing KV data, lightweight analog of the Redis

MAKEFILE:

    build   - Compiles the binaries for server and for client and creates the ready-to-use docker container (use make run)
    test    - Runs Unit and Integration tests for the server. File coverage.out will apeear in the serv dir
    check   - Runs subsequently "go vet", "goimports", "golint" on the project. Fails if any errors occur.
    run     - Starts the server and the client with default configuration (port: 9090, storage: RAM)

serv:
    -p, --port  - The port for listening on (default 9090)
    -m, --mode  - The possible storage option
            <memory> - use RAM ad a storage (default)
            <disk> - use disk as a storage
commands:
    SET <key> <value>   - updates one key at a time with the given value.
    GET <key>           - returns tuple of the value and the key state.
    DEL <key>           - removes one key at a time and returns deleted value.


client:
    -p, --port  - The port for listening on (default 9090)
    -h, --host  - The host to connect to the server (default: 127.0.0.1)
