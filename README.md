# Prerequisites
Docker Desktop

Init the project
`go mod init evaluation-sys-kafka`

Downloa
`go get github.com/confluentinc/confluent-kafka-go/kafka`

Install Confluent CLI
`brew install confluentinc/tap/cli`

Start the Kafka broker
`confluent local kafka start`
