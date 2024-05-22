#!/bin/bash

# Spin up the container
# Generate a unique container name with a timestamp
new_container_name="couple-bot-service-$(date +%s)"

echo "New container name: $new_container_name"

sudo docker run --rm -d --name "$new_container_name" \
    -e TELEGRAM_BOT_TOKEN="$2" \
    gcr.io/linen-shape-420522/couple-bot-service:latest

# Check if the new container is up and running
while [ "$(sudo docker inspect -f '{{.State.Running}}' "$new_container_name")" != "true" ]; do
    echo "Waiting for container to start: $new_container_name"
    sleep 1
done

echo "New container is running: $new_container_name"
# Find and stop the old container
old_container=$(sudo docker ps --filter "name=couple-bot-service" --format "{{.Names}}" | grep -v "$new_container_name")
if [ ! -z "$old_container" ]; then
    sudo docker stop "$old_container"
    echo "Old container stopped: $old_container"
fi

