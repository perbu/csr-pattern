#!/usr/bin/env bash

set -e

# Set the API endpoint
URL="http://localhost:8080/mykey" # Replace with your actual API endpoint
VALUE1="myvalue1"
VALUE2="myvalue2"

# Insert the key-value pair
echo "Inserting key-value pair..."
http --json POST "$URL"  value="$VALUE1"

# Get the value
echo "Getting value..."
http --json GET "$URL"

# Update the value
echo "Updating value..."
http --json PUT "$URL" value="$VALUE2"

# Get the updated value
echo "Getting updated value..."
http --json GET "$URL"

# Delete the key-value pair
echo "Deleting key-value pair..."
http DELETE "$URL"

echo "Done."