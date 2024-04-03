package controllers

import (
	kafkaUtils "evaluation-sys-kafka/internal/kafka"
	courseModels "evaluation-sys-kafka/pkg/courses/models"
	"evaluation-sys-kafka/pkg/users/models"
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Producer       *kafka.Producer
	ConsumerOutput []interface{}
	Students       []models.Student
	Professors     []models.Professor
}

// TODO: update based on action_type
func (c *Controller) GetCourses(context *gin.Context) {
	var courses []courseModels.Course // Initialize a slice to hold Course instances

	for index, element := range c.ConsumerOutput {
		//fmt.Printf("Index: %d, Value: %v, Type: %T\n", index, element, element)
		if courseMap, ok := element.(map[string]interface{}); ok {
			course := courseModels.Course{
				ID:   fmt.Sprint(courseMap["id"]),
				Name: fmt.Sprint(courseMap["name"]),
			}
			courses = append(courses, course)
		} else {
			fmt.Printf("Error: element at index %d cannot be converted to Course\n", index)
		}
	}
	context.JSON(http.StatusOK, gin.H{"courses": courses})
}

// ASSUMPTION: No authentication; student record is created with a simple POST request
func (c *Controller) CreateStudent(context *gin.Context) {
	var request models.Student
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, student := range c.Students {
		if student.ID == request.ID {
			context.JSON(http.StatusConflict, gin.H{"error": "Student with such ID already present"})
			return
		}
	}
	c.Students = append(c.Students, request)
	fmt.Println("(CreateStudent) > In-memory Students: ", c.Students)
	err := kafkaUtils.ProduceMessage(c.Producer, "add", "student", request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// wait for all messages to be acknowledged
	c.Producer.Flush(15 * 1000)
	context.JSON(http.StatusCreated, gin.H{"message": "Student created successfully"})
}

func (c *Controller) CreateProfessor(context *gin.Context) {
	var request models.Professor
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, professor := range c.Professors {
		if professor.ID == request.ID {
			context.JSON(http.StatusConflict, gin.H{"error": "Professor with such ID already present"})
			return
		}
	}
	c.Professors = append(c.Professors, request)
	fmt.Println("(CreateProfessor) > In-memory Professors: ", c.Professors)
	// TODO: here you could produce some events, for now it's good as it is
	context.JSON(http.StatusCreated, gin.H{"message": "Professor created successfully"})
}
