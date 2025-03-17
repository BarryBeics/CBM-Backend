#!/bin/bash

echo "Running priceData at $(date)" >> /var/log/priceData.log
./usr/local/bin/microservice-binaries/priceData --server=nats://nats:4222 --loglevel=debug --logcolor >> /var/log/priceData.log 2>&1
echo "Finished priceData at $(date) with exit code $?" >> /var/log/priceData.log

# Keep the container alive
tail -f
