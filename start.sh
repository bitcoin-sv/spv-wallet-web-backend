#!/usr/bin/env bash

# Ensure .local directory exists
mkdir -p .local

# Download newest version of the start.sh and docker-compose.yml
echo ""
echo "Downloading newest version of start.sh from bitcoin-sv/spv-wallet repository:"
curl -o .local/start.sh https://raw.githubusercontent.com/bitcoin-sv/spv-wallet/main/start.sh
echo ""
echo "Downloading newest version of docker-compose.yml from bitcoin-sv/spv-wallet repository:"
curl -o .local/docker-compose.yml https://raw.githubusercontent.com/bitcoin-sv/spv-wallet/main/docker-compose.yml

# ensure that script is runnable
chmod +x .local/start.sh

# run the start.sh in context of .local directory
cd .local && ./start.sh "$@"

# react on CTRL+C
trap 'cd .local && docker compose down' SIGINT
