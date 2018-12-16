################ Preparation ####################
REPOSITORY_PATH_SERVER := kv_storage_server
REPOSITORY_PATH_CLIENT := kv_storage_client
REGISTRY      := agri88
VERSION       := $(shell cat VERSION)
BECOME        := sudo -E
DOCKER_IMAGE   := $(REGISTRY)/kv_storage:$(VERSION)
DOCKER_BUILDER := golang:1.11 
DOCKER_RUNNER_SERVER =  docker run --rm -v $(CURDIR)/server/:/go/src/$(REPOSITORY_PATH_SERVER)
DOCKER_RUNNER_SERVER += $(DOCKER_ENVS) -w /go/src/$(REPOSITORY_PATH_SERVER)
DOCKER_RUNNER_CLIENT =  docker run --rm -v $(CURDIR)/client/:/go/src/$(REPOSITORY_PATH_CLIENT)
DOCKER_RUNNER_CLIENT += $(DOCKER_ENVS) -w /go/src/$(REPOSITORY_PATH_CLIENT)
 # search files for fmt and check targets, excluding "vendor" folder
SEARCH_GOFILES =  find -not -path '*/vendor/*' -type f -name "*.go"
BUILDER_SERVER       = $(DOCKER_RUNNER_SERVER) $(DOCKER_BUILDER)
BUILDER_CLIENT       = $(DOCKER_RUNNER_CLIENT) $(DOCKER_BUILDER)
PORT = 9090
CONTAINER_NAME = kv_server
IMAGE_ID := $(shell $(BECOME) docker images -q $(DOCKER_IMAGE))
 ################ End Preparation ####################
 ################ Binary Target ####################
build: clean server_bin client_bin

.PHONY: server_bin
server_bin:
	mkdir ./server/common_files && \
	mkdir ./server/log && \
	$(BECOME) $(BUILDER_SERVER) go build -o $@ ./$(@D)
	$(BECOME) docker build -t $(DOCKER_IMAGE) ./
.PHONY: client_bin
client_bin:
	$(BECOME) $(BUILDER_CLIENT) go build -o $@ ./$(@D)
 ################ Clean Target ####################
.PHONY: clean
clean:
	$(BECOME) $(RM) ./server/server_bin
	$(BECOME) $(RM) ./client/client_bin 
	$(BECOME) $(RM) -r ./server/log
	$(BECOME) $(RM) -r ./server/common_files
ifneq ($(shell $(BECOME) docker ps -q -f name=$(CONTAINER_NAME)), )
	$(BECOME) docker rm -f $(CONTAINER_NAME)	
endif
ifneq ($(IMAGE_ID), )
	$(BECOME) docker image rm $(IMAGE_ID)
endif
 ################ Format and Validate Targets ####################
.PHONY: check
check:
	$(BECOME) docker run --name $(CONTAINER_NAME)_check -d $(DOCKER_IMAGE)
	$(BECOME) docker exec -it $(CONTAINER_NAME)_check gofmt -s -l ./server 
	$(BECOME) docker exec -it $(CONTAINER_NAME)_check golint ./server 
	$(BECOME) docker exec -it $(CONTAINER_NAME)_check go vet ./server 
	$(BECOME) docker exec -it $(CONTAINER_NAME)_check gofmt -s -l ./client  
	$(BECOME) docker exec -it $(CONTAINER_NAME)_check golint ./client
	$(BECOME) docker exec -it $(CONTAINER_NAME)_check go vet ./client 
	$(BECOME) docker rm -f $(CONTAINER_NAME)_check 
 ################ Docker Targets ####################
.PHONY: run
ifneq ($(shell $(BECOME) docker ps -a -q -f name= $(CONTAINER_ID)), )
run:
	echo container is created
else
ifneq ($(IMAGE_ID), )
run:
	$(BECOME) docker run --name $(CONTAINER_NAME) -p $(PORT):9090 -v $(CURDIR)/server/log/:/go/src/log/ -v $(CURDIR)/server/common_files/:/go/src/common_files/ -d $(DOCKER_IMAGE)
else
run: build 
	$(BECOME) docker run --name $(CONTAINER_NAME) -p $(PORT):9090 -v $(CURDIR)/server/log/:/go/src/log/ -v $(CURDIR)/server/common_files/:/go/src/common_files/ -d $(DOCKER_IMAGE)
endif

endif

.PHONY: stop
ifneq ($(shell $(BECOME) docker ps -a -q -f name= $(CONTAINER_ID)), )
stop:
	$(BECOME) docker rm -f $(CONTAINER_NAME)
else
stop:
	echo container not exist
endif
 ################ Test Targets ####################
.PHONY: test
test:	
	$(BECOME) docker exec -it $(CONTAINER_NAME) go test ./server/ -coverprofile ./common_files/cover.out
