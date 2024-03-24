package controllers

import (
	"evaluation-sys-kafka/pkg/courses/models"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Producer *kafka.Producer
}

func (c *Controller) CreateCourse(context *gin.Context) {
	var request models.Course
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Convert the course model to a map[string]string
	/*courseMap := map[string]string{
		"course_id":   request.ID,
		"course_name": request.Name,
	}*/

	topic := "course"
	err := c.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte("trykey"),
		Value:          []byte("tryval"),
	}, nil)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: understand whether this is necessary
	//c.Producer.Flush(15 * 1000) // Wait for 15 seconds for the producer to finish

	context.JSON(http.StatusOK, gin.H{"message": "Course created successfully"})
}
