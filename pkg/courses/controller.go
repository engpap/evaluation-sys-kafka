package courses

import "evaluation-sys-kafka/internal/kafka"

func Run() {
	kafka.CreateProducer("course", map[string]string{"course_id": "23"})
}
