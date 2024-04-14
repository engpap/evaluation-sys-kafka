docker build -f cmd/user-service/Dockerfile -t user-service .
docker run -p 8083:8083 user-service