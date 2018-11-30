################ Preparation ####################
REPOSITORY_PATH := ./TinyRedisServer

VERSION       := $(shell cat VERSION)

BECOME        := sudo -E
DOCKER_IMAGE   := $(REGISTRY)/TinyRedisServer:$(VERSION)
DOCKER_BUILDER := golang:1.11
DOCKER_RUNNER =  docker run --rm -v $(CURDIR):/go/src/$(REPOSITORY_PATH)
DOCKER_RUNNER += $(DOCKER_ENVS) -w /go/src/$(REPOSITORY_PATH)

# search files for fmt and check targes, excluding "vendor" folder
SEARCH_GOFILES =  find -not -path '*/vendor/*' -type f -name "*.go"
BUILDER       = $(DOCKER_RUNNER) $(DOCKER_BUILDER)

################ End Preparation ####################



################ Binary Target ####################

TinyRedisServer: VERSION check
	$(BECOME) $(BUILDER) go build -o $@ ./$(@D)

################ Clean Target ####################
.PHONY: clean
clean:
	$(BECOME) $(RM) TinyRedisServer

################ Format and Validate Targets ####################
.PHONY: fmt
fmt:
	$(BECOME) $(BUILDER) $(SEARCH_GOFILES) -exec gofmt -s -w {} \;

.PHONY: check
check:
	$(BECOME) $(BUILDER) sh -xc '\
		test -z "`$(SEARCH_GOFILES) -exec gofmt -s -l {} \;`" \
		&& test -z "`$(SEARCH_GOFILES) -exec golint {} \;`"'

################ Docker Targets ####################

################ Test Targets ####################

