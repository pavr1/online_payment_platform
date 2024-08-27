#!/bin/bash

# Set the environment variables
VARIABLES=(
  "AUTH_PORT=8081"
)

echo "Setting environment variables in namespace $1"

# Loop through the variables and add them to the environment
for variable in "${VARIABLES[@]}"; do
  kubectl -n $1 set env deployment/auth ${variable}
done

echo "Done!"