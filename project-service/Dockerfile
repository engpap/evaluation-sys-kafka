# INSTRUCTIONS:
# AWS expects the Dockerfile to be in the root of the project => copy temporarily to the root
# cp project-service/Dockerfile .
#
# Build the Docker image using the following command  from the root of the project
# docker build -f Dockerfile -t project-service .
# 
# Remove the temporary Dockerfile from the root
# rm Dockerfile

#
# Run the Docker container using the following command
# docker run -p 8080:8080 project-service

FROM golang:1.21.5

# Set the working directory inside the container
WORKDIR /app

# Copy only what's needed into the container
COPY project-service/ /app/project-service/
COPY client.properties /app/client.properties
#COPY go.work /app/go.work
#COPY go.work.sum /app/go.work.sum
COPY Dockerfile /app/Dockerfile


WORKDIR /app/project-service
RUN go mod tidy

# (DEBUG): List the files in the current directory
# Open Docker app and see `logs` section
# Uncomment the following line to enable directory listing in logs
RUN ls -la 
RUN pwd
RUN ls -ls /app
RUN ls -ls /app/project-service
RUN ls -ls /app/project-service/cmd

# Build the Go app
# Current directory is /app; here we can access the cmd/app directory
# The -o flag specifies the output path of the executable
RUN go build -o /app/main /app/project-service/cmd/app

# Expose port for the application
EXPOSE 8080

# Run the outyet command by default when the container starts.
# Run the main executable with the client.properties file as an argument
# The executable is in the current directory
# The client.properties file is in the /app directory, which is 2 levels up
ENTRYPOINT ["/app/main", "/app/client.properties"]