#!/bin/bash
set -e


echo serv go vet
go vet
echo serv golint
golint
echo serv goimports
goimports
