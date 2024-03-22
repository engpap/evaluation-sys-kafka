package users

import "evaluation-sys-kafka/internal/kafka"

func Run() {
	// TEST, CHANGE WITH DIFFERENT PARAMETER
	kafka.CreateConsumer("course")
}
