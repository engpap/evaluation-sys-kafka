# Prerequisites
Docker Desktop

# Instructions
Init the project
`go mod init evaluation-sys-kafka`

Download
`go get github.com/confluentinc/confluent-kafka-go/kafka`

Install Confluent CLI
`brew install confluentinc/tap/cli`

Start the Kafka broker
`confluent local kafka start`

Note the Plaintext Ports printed in your terminal; paste the following configuration data into a file named getting-started.properties, substituting the plaintext port(s) output.
`bootstrap.servers=localhost:<PLAINTEXT PORTS>`

Create a topic
`confluent local kafka topic create <topic_name>`

Compile producer (assuming course-service is the producer)
`go build -o out/producer ./cmd/course-service`

Compile consumer (assuming user-service is the producer)
`go build -o out/consumer ./cmd/user-service`

Run producer
`./out/producer getting-started.properties`

Run consumer
`./out/consumer getting-started.properties`

Type Ctrl-C to terminate consumer application

Shut down Kafka when you are done with it
`confluent local kafka stop`