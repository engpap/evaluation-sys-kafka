package controllers

import (
	kafkaUtils "evaluation-sys-kafka/internal/kafka"
	"evaluation-sys-kafka/pkg/courses/models"
	usersModels "evaluation-sys-kafka/pkg/users/models"
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Producer *kafka.Producer
	// In-memory data structures
	Courses     []models.Course
	Enrollments []models.Enrollment
	// In-memory data structures populated by consumers
	Students []usersModels.Student
}

func (c *Controller) GetCourses(context *gin.Context) {
	context.JSON(http.StatusOK, c.Courses)
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
	context.JSON(http.StatusCreated, gin.H{"message": "Course created successfully"})
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

// checks if enrollment already exists in memory. if not, it checks existance of course and student provided.
// students needs to be fetched from kafka topic (user-service producers pushes them into `student` topic)
func (c *Controller) EnrollStudentInCourse(context *gin.Context) {
	var request models.Enrollment
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// prevent adding when already present
	fmt.Println("(EnrollStudentInCourse) > In-memory Enrollments: ", c.Enrollments)
	for _, enrollment := range c.Enrollments {
		if request.StudentID == enrollment.StudentID && request.CourseID == enrollment.CourseID {
			context.JSON(http.StatusConflict, gin.H{"error": "Student already enrolled in this course"})
			return
		}
	}
	// check course exists
	fmt.Println("(EnrollStudentInCourse) > In-memory Courses: ", c.Courses)
	found := false
	for _, course := range c.Courses {
		if request.CourseID == course.ID {
			found = true
			break
		}
	}
	if !found {
		context.JSON(http.StatusNotFound, gin.H{"error": "Bad request. You're trying to enroll in a course that does not exists."})
		return
	}
	found = false
	for _, student := range c.Students {
		if request.StudentID == student.ID {
			found = true
			break
		}
	}
	if !found {
		context.JSON(http.StatusNotFound, gin.H{"error": "Bad request. You're trying to enroll a student that does not exists."})
		return
	}
	// at this point enrollment is valid => update in-memory state
	c.Enrollments = append(c.Enrollments, request)
	fmt.Println("(EnrollStudentInCourse) > In-memory Enrollments: ", c.Enrollments)
	err := kafkaUtils.ProduceMessage(c.Producer, "add", "enrollment", request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// wait for all messages to be acknowledged
	c.Producer.Flush(15 * 1000)
	context.JSON(http.StatusCreated, gin.H{"message": "Enrollment created successfully"})
}

// / CALLBACK FUNCTIONS

func (c *Controller) UpdateStudentInMemory(action_type string, data interface{}) {
	if action_type == "add" {
		c.saveStudentInMemory(data)
	} else {
		fmt.Printf("Error: action type %s not supported\n", action_type)
	}
}
func (c *Controller) saveStudentInMemory(data interface{}) {
	if studentMap, ok := data.(map[string]interface{}); ok {
		student := usersModels.Student{
			ID: fmt.Sprint(studentMap["id"]),
		}
		c.Students = append(c.Students, student)
		fmt.Println("In-memory Students: ", c.Students)
	} else {
		fmt.Printf("Error: data cannot be converted to Student\n")
	}

}
