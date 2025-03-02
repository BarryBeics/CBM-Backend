#!/bin/bash
set -x

# Ensure /usr/local/bin is in PATH
export PATH=$PATH:/usr/local/bin

echo "Starting Nex and registering binaries with NEX..."

# Loop through all binaries in the microservice-binaries directory
for binary in /usr/local/bin/microservice-binaries/*; do
   if [[ -x "$binary" && -f "$binary" ]]; then
       binary_name=$(basename "$binary")
       echo "Attempting to register $binary_name binary with Nex..."
       "$binary" --server=nats://nats:4222 --loglevel=debug --logcolor || echo "Failed to register $binary_name"
   else
       echo "Skipping $binary - not an executable file."
   fi
done

# Keep the container running to allow Nex operations
tail -f /dev/null