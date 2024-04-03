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
	StudentConsumerOutput []interface{}
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
// TODO: should i fetch only new students to speed things up (?) or create a listener instead of fetching on request
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
	// check student exists (through consuming)
	fmt.Println("(EnrollStudentInCourse) > Students to consume: ", c.StudentConsumerOutput)
	// iterates over the StudentConsumerOutput to fetch the students
	var students []usersModels.Student
	for index, element := range c.StudentConsumerOutput {
		if studentMap, ok := element.(map[string]interface{}); ok {
			student := usersModels.Student{
				ID: fmt.Sprint(studentMap["id"]),
			}
			students = append(students, student)
		} else {
			fmt.Printf("Error: element at index %d cannot be converted to Student\n", index)
		}
	}

	fmt.Println("(EnrollStudentInCourse) > In-memory Students: ", students)
	fmt.Println("WARNING > (EnrollStudentInCourse) > Implement a listener instead for students")
	// check student exists
	found = false
	for _, student := range students {
		if student.ID == request.StudentID {
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
