# Prerequisites
- Docker Desktop
- Confluent CLI

# How to run the backend microservices
Make sure to set the boolean variable debug in the services to true if you want to run the project locally. If you want to run the project on AWS, set the variable to false. You can find this variable in the transport.go file of each service. If debug is set to true and the services are running on AWS, they won't be reachable from the CLI frontend because they will run on a port different from 8080, which is not what AWS expects. As a result, attempting to connect to them would lead to a 502 Bad Gateway error.
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
go run cli-app/cmd/app/main.go
```

# Observation
The package [kafka-wrapper-go](https://github.com/engpap/kafka-wrapper-go) is part of this project. It is a wrapper around the confluent-kafka-go library. It is used to simplify the Kafka producer and consumer applications. The wrapper provides a simple API to interact with Kafka topics and provides a fault recovery functionality. The wrapper is used in the producer and consumer applications of the project. 