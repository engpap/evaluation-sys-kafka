# This script zips the services for deployment on AWS
# Store them in the `out` directory

# For each microservice:
#   - copy the Dockerfile to the root directory
#   - zip the microservice
#   - remove the Dockerfile from the root directory
# This is done so that AWS can find the Dockerfile in the source bundle; otherwise, AWS will not be able to build the Docker image.

# pkg is included in every microservice because it contains the shared code
# client.properties is included in every microservice because it contains the Kafka broker configuration
# go.work and go.work.sum are included in every microservice because they contain the Go module dependencies between the microservices

cp user-service/Dockerfile ./
zip -r out/user-service.zip user-service pkg client.properties go.work go.work.sum Dockerfile
rm Dockerfile

cp project-service/Dockerfile ./
zip -r out/project-service.zip project-service pkg client.properties go.work go.work.sum Dockerfile
rm Dockerfile

cp registration-service/Dockerfile ./
zip -r out/registration-service.zip registration-service pkg client.properties go.work go.work.sum Dockerfile
rm Dockerfile

cp course-service/Dockerfile ./
zip -r out/course-service.zip course-service pkg client.properties go.work go.work.sum Dockerfile
rm Dockerfile