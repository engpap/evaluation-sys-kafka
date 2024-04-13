#!/bin/bash

# This script is used to start the Kafka broker, create Kafka topics, compile the micro-services, and start the micro-services in new terminals.

# Function to stop Kafka and perform cleanup
cleanup() {
  echo "Stopping Kafka..."
  confluent local kafka stop
  echo "Kafka stopped."
  exit 0
}
# Trap SIGINT (Ctrl+C) and call the cleanup function
trap cleanup SIGINT

# Exit script on error. Meaning: If any command fails, the script will stop executing
set -e

# Stop Kafka broker if it is running
confluent local kafka stop 

# Start Kafka broker
confluent local kafka start

# Ask for the Kafka broker port and update getting-started.properties
read -p "Enter Kafka Plaintext Ports: " kafka_port
echo "bootstrap.servers=localhost:$kafka_port" > getting-started.properties

echo "Creating Kafka topics..."
declare -a topic_names=("course" "student" "project" "submission" "grade" "enrollment")
for topic_name in "${topic_names[@]}"
do
  confluent local kafka topic create $topic_name
done

echo "Compiling micro-services..."
go build -o out/user ./cmd/user-service
go build -o out/project ./cmd/project-service
go build -o out/course ./cmd/course-service
go build -o out/registration ./cmd/registration-service

echo "Starting micro-services in new terminals..."
# Check for OS and launch terminals accordingly
case "$(uname)" in "Linux")
    xterm -e "./out/user getting-started.properties" &
    xterm -e "./out/project getting-started.properties" &
    xterm -e "./out/course getting-started.properties" &
;;
"Darwin")
    osascript -e 'tell application "Terminal" to do script "cd '$(pwd)' && ./out/user getting-started.properties"'
    osascript -e 'tell application "Terminal" to do script "cd '$(pwd)' && ./out/project getting-started.properties"'
    osascript -e 'tell application "Terminal" to do script "cxd '$(pwd)' && ./out/course getting-started.properties"'
    osascript -e 'tell application "Terminal" to do script "cd '$(pwd)' && ./out/registration getting-started.properties"'
  ;;
*)
  echo "Unsupported OS for launching new terminals."
  ;;
esac

# Wait indefinitely until SIGINT is received
echo "Press Ctrl+C to stop Kafka and exit."
while true; do sleep 1; done