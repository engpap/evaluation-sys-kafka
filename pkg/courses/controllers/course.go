package controllers

import (
	"encoding/json"
	"evaluation-sys-kafka/pkg/courses/models"
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Producer *kafka.Producer
	Courses  []models.Course
}

func (c *Controller) CreateCourse(context *gin.Context) {
	var request models.Course
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update state
	c.Courses = append(c.Courses, request)
	fmt.Println("In-memory Courses: ", c.Courses)
	// Convert into a byte array
	data, err := json.Marshal(request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Send the message to the Kafka topic
	topic := "course"
	err = c.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte("data"),
		Value:          data,
	}, nil)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Waiit for all messages to be acknowledged
	c.Producer.Flush(15 * 1000)
	context.JSON(http.StatusOK, gin.H{"message": "Course created successfully"})
}
