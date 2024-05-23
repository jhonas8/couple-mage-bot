#!/bin/bash

# Spin up the container
# Generate a unique container name with a timestamp
new_container_name="couple-bot-service-$(date +%s)"

echo "New container name: $new_container_name"

if [ -z "$1" ]; then
    if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
        echo "No telegram bot token provided as first argument or in environment."
        exit 1
    fi
else
    telegram_bot_token="$1"
fi

if [ -z "$2" ]; then
    github_ssh_token=""
else
    github_ssh_token="$2"
fi

# Request the value of the secret FIREBASE_SECRET_PAAL_BASE64 from secrets manager
# firebase_secret_base64=$(gcloud auth activate-service-account --key-file=$gcp_service_account_base64 && gcloud secrets versions access latest --secret=FIREBASE_SECRET_PAAL_BASE64 --project=crispai)


# Pull the latest changes from the repository
if [ -d "couple-mage-bot" ]; then
    sudo rm -rf couple-mage-bot
fi

# set a git username and password
sudo git config --global user.name "jhonas8"
sudo git config --global user.email "love.android62@gmail.com"

# git credentials cache

sudo git clone "https://$github_ssh_token@github.com/jhonas8/couple-mage-bot.git"
cd couple-mage-bot
# Build the Docker image from the Dockerfile in the current directory
sudo docker build -t "$new_container_name" .

# Run the newly built Docker image
sudo docker run --rm -d --name "$new_container_name" \
    -e TELEGRAM_BOT_TOKEN="$telegram_bot_token" \
    "$new_container_name"

# Check if the new container is up and running
while [ "$(sudo docker inspect -f '{{.State.Running}}' "$new_container_name")" != "true" ]; do
    echo "Waiting for container to start: $new_container_name"
    sleep 1
done

echo "New container is running: $new_container_name"
# Find, stop, and remove the old container along with its image
old_container=$(sudo docker ps --filter "name=couple-bot-service" --format "{{.Names}}" | grep -v "$new_container_name")
if [ ! -z "$old_container" ]; then
    old_image=$(sudo docker inspect --format='{{.Image}}' "$old_container")
    sudo docker stop "$old_container"
    echo "Old container stopped: $old_container"
    sudo docker rm "$old_container"
    echo "Old container removed: $old_container"
    sudo docker rmi "$old_image"
    echo "Old container image removed: $old_image"
fi
