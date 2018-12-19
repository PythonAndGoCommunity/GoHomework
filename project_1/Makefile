build:
	cd client_source; go build client.go; cd ..;
	cd server_source; go build server.go; cd ..;
	mv -f client_source/client .
	mv -f server_source/server .
	sudo docker build -f DockerFile -t custom_redis .

test:
	cd client_source; go test -cover -coverprofile=coverage_client.out; cd ..;
	cd server_source; go test -cover -coverprofile=coverage_server.out; cd ..;
	mv -f client_source/coverage_client.out .
	mv -f server_source/coverage_server.out .

check:
	cd server_source; go vet; golint; goimports .; cd ..;
	cd client_source; go vet; golint; goimports .; cd ..;
	cd commonVariables; go vet; golint; goimports .; cd ..;

run:
	sudo docker run --rm --name my_redis -it custom_redis
