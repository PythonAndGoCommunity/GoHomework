#!/bin/bash
set -e

echo starting server
./serv/serv & &>/dev/nul
echo 'waiting for cli to start'
sleep 2
echo starting client
./client/client
