#!/usr/bin/env bash

# run chains
printf "running chains.."
{
    docker-compose -f ./e2e/evm-evm/docker-compose.e2e.yml up -d
} || {
	exit
}

printf "deploying local environment.."
# run local-setup
{
    go run e2e/evm-evm/example/main.go local-setup
} || {
	exit
}
