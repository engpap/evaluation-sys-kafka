package controllers

import (
	"eval-sys-user-service/cmd/internal/models"
	"fmt"
	"net/http"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	kafkaWrapper "github.com/engpap/kafka-wrapper-go/pkg"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	mu       sync.Mutex
	Producer *kafka.Producer
	// In-memory data structures that will be populated through consuming
	Students   []models.Student
	Professors []models.Professor
}

// ASSUMPTION: No authentication; student record is created with a simple POST request
func (c *Controller) CreateStudent(context *gin.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
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
	c.mu.Lock()
	defer c.mu.Unlock()
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
	err := kafkaWrapper.ProduceMessage(c.Producer, "add", "professor", request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Professor created successfully"})
}

// CALLBACK FUNCTIONS

func (c *Controller) UpdateStudentInMemory(action_type string, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if action_type == "add" {
		c.saveStudentInMemory(data)
	} else {
		fmt.Printf("Invalid action type: %s\n", action_type)
	}
}

func (c *Controller) saveStudentInMemory(data interface{}) {
	if studentMap, ok := data.(map[string]interface{}); ok {
		student := models.Student{
			ID: fmt.Sprint(studentMap["id"]),
		}
		c.Students = append(c.Students, student)
		fmt.Println("In-memory Students: ", c.Students)
	} else {
		fmt.Printf("Error: data cannot be converted to Student\n")
	}
}

func (c *Controller) UpdateProfessorInMemory(action_type string, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if action_type == "add" {
		c.saveProfessorInMemory(data)
	} else {
		fmt.Printf("Invalid action type: %s\n", action_type)
	}
}

func (c *Controller) saveProfessorInMemory(data interface{}) {
	if professorMap, ok := data.(map[string]interface{}); ok {
		professor := models.Professor{
			ID: fmt.Sprint(professorMap["id"]),
		}
		c.Professors = append(c.Professors, professor)
		fmt.Println("In-memory Professors: ", c.Professors)
	} else {
		fmt.Printf("Error: data cannot be converted to Professor\n")
	}
}
