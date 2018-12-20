#!/bin/bash
set -e


echo client go vet
go vet
echo client golint
golint
echo client goimports
goimports
