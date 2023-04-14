#!/bin/bash

# Stop all running containers
docker stop $(docker ps -aq)

# Remove all stopped containers
docker rm -f $(docker ps -aq)

# Remove all images
docker rmi -f $(docker images -q)

rm ./example/log_evaluate.txt
rm ./example/scripts/all.txt
