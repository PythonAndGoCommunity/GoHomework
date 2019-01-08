#=============#
#    Redis    #
# ============#

all: get install

fmt:
	@ go fmt ./...

get:
	@ govendor sync

install:
	@ go install ./...

test:
	@ go test -p 1 ./...
