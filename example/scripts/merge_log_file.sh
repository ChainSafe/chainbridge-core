#!/bin/bash

# Define array of container names
containers=("relayer1" "relayer2" "relayer3")

# Loop through containers
for container in "${containers[@]}"
do
  # Enter the container
  docker exec -it $container bash -c "cat ./log_evaluate.txt" >> all.txt
done
