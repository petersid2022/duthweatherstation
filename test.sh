#!/bin/bash

# Function to get the timestamp in the desired format
getFormattedTime() {
  local date_input=$1
  if [ -z "$date_input" ]; then
    date -u +"%Y-%m-%dT%H:%M:%SZ"
  else
    date -u -d "$date_input" +"%Y-%m-%dT%H:%M:%SZ"
  fi
}

# Function to create a JSON payload
createPayload() {
  local sensor=$1
  local timestamp=$2
  local data=$3
  echo "{\"SensorName\":\"$sensor\",\"SensorData\":{\"timestamp\":\"$timestamp\",$data}}"
}

# Function to send data
sendData() {
  local ad_value=$1
  local h=$2
  local t=$3
  local p=$4
  local date_input=$5

  local serverName="https://duthweather.azurewebsites.net/api/add" # Replace with your server name
  local timestamp=$(getFormattedTime "$date_input")

  # Create JSON payloads
  local payload_dht11=$(createPayload "dht11" "$timestamp" "\"temperature\":$t,\"humidity\":$h")
  local payload_bmp180=$(createPayload "bmp180" "$timestamp" "\"pressure\":$p")
  local payload_mq135=$(createPayload "mq135" "$timestamp" "\"gas_level\":$ad_value")

  echo "Uploading data"

  # Send POST requests using curl
  curl -X POST -H "Content-Type: application/json" -d "$payload_dht11" "$serverName"
  curl -X POST -H "Content-Type: application/json" -d "$payload_bmp180" "$serverName"
  curl -X POST -H "Content-Type: application/json" -d "$payload_mq135" "$serverName"

  echo "Done!"
}

# Check if the correct number of arguments is provided
if [ "$#" -lt 4 ] || [ "$#" -gt 5 ]; then
  echo "Usage: $0 ad_value humidity temperature pressure [timestamp]"
  echo "If no timestamp is provided, the current time will be used."
  exit 1
fi

# Call the sendData function with the provided arguments
sendData "$1" "$2" "$3" "$4" "$5"
