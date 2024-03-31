package controllers

import (
	kafkaUtils "evaluation-sys-kafka/internal/kafka"
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
	// prevent adding when id already present
	for _, course := range c.Courses {
		if request.ID == course.ID {
			context.JSON(http.StatusConflict, gin.H{"error": "Course with such ID already present"})
			return
		}
	}
	// update in-memory state
	c.Courses = append(c.Courses, request)
	fmt.Println("(CreateCourse) > In-memory Courses: ", c.Courses)
	err := kafkaUtils.ProduceMessage(c.Producer, "add", "course", request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// wait for all messages to be acknowledged
	c.Producer.Flush(15 * 1000)
	context.JSON(http.StatusOK, gin.H{"message": "Course created successfully"})
}

func (c *Controller) DeleteCourse(context *gin.Context) {
	courseID := context.Param("course-id")
	if courseID == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Course ID to delete not provided"})
		return
	}
	// Iterate over courses global variable and delete when there's a match
	for index, course := range c.Courses {
		if course.ID == courseID {
			c.Courses = append(c.Courses[:index], c.Courses[index+1:]...)
			fmt.Println("(DeleteCourse) > In-memory Courses: ", c.Courses)
			err := kafkaUtils.ProduceMessage(c.Producer, "delete", "course", course)
			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			context.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
			return
		}
	}
	context.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
}
