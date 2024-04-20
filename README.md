# Prerequisites
- Docker Desktop
- Confluent CLI

# How to run the backend microservices
## Locally
Run:
```
./scripts/run-local.sh
```
This script will start the Kafka broker and create the necessary topics. It will also start the producer and consumer applications by launching 4 microservices on separate terminals.

## Locally on Docker
Run:
```
docker build -f cmd/course-service/Dockerfile -t course-service .
docker run -p 8090:8090 course-service
```
```
docker build -f cmd/project-service/Dockerfile -t project-service .
docker run -p 8091:8091 project-service
```
```
docker build -f cmd/register-service/Dockerfile -t register-service .
docker run -p 8092:8092 register-service
```
```
docker build -f cmd/user-service/Dockerfile -t user-service .
docker run -p 8093:8093 user-service
```

## Remotely on AWS
### Deployment on AWS
Run:
```
./scripts/zip-services.sh
```
This script will zip the services with their dependencies. Login to your AWS account and upload the zip files to the Elastic Beanstalk environment.

# How to run the CLI frontend
Make sure to set the boolean variable `debug` in the CLI app to `true` if you want to run the project locally. If you want to run the project on AWS, set the variable to `false`.
Then run:
```
go run cmd/cli-frontend/main.go
```
