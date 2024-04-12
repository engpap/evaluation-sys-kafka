#!/bin/bash

# This script is used to compile and run the microservices with the condiguration to connect them to the Kafka broker running on the cloud.

echo "Compiling micro-services..."
go build -o out/user ./cmd/user-service
go build -o out/project ./cmd/project-service
go build -o out/course ./cmd/course-service
go build -o out/registration ./cmd/registration-service

echo "Starting micro-services in new terminals..."
# Check for OS and launch terminals accordingly
case "$(uname)" in "Linux")
    xterm -e "./out/user client.properties" &
    xterm -e "./out/project client.properties" &
    xterm -e "./out/course client.properties" &
;;
"Darwin")
    osascript -e 'tell application "Terminal" to do script "cd '$(pwd)' && ./out/user client.properties"'
    osascript -e 'tell application "Terminal" to do script "cd '$(pwd)' && ./out/project client.properties"'
    osascript -e 'tell application "Terminal" to do script "cxd '$(pwd)' && ./out/course client.properties"'
    osascript -e 'tell application "Terminal" to do script "cd '$(pwd)' && ./out/registration client.properties"'
  ;;
*)
  echo "Unsupported OS for launching new terminals."
  ;;
esac

# Wait indefinitely until SIGINT is received
echo "Press Ctrl+C to stop Kafka and exit."
while true; do sleep 1; done