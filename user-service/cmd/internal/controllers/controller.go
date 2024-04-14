package controllers

import (
	"eval-sys-user-service/cmd/internal/models"
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	kafkaWrapper "github.com/engpap/kafka-wrapper-go/pkg"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Producer   *kafka.Producer
	Students   []models.Student
	Professors []models.Professor
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
	err := kafkaWrapper.ProduceMessage(c.Producer, "add", "student", request)
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
