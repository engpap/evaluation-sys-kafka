package controllers

import (
	kafkaUtils "evaluation-sys-kafka/internal/kafka"
	"evaluation-sys-kafka/pkg/projects/models"
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Producer    *kafka.Producer
	Projects    []models.Project
	Submissions []models.Submission
	Grades      []models.Grade
}

func (c *Controller) CreateProject(context *gin.Context) {
	var request models.Project
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, project := range c.Projects {
		if project.ID == request.ID {
			context.JSON(http.StatusConflict, gin.H{"error": "Project with such ID already present"})
			return
		}
	}
	c.Projects = append(c.Projects, request)
	fmt.Println("(CreateProject) > In-memory Projects: ", c.Projects)
	err := kafkaUtils.ProduceMessage(c.Producer, "add", "project", request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// wait for all messages to be acknowledged
	c.Producer.Flush(15 * 1000)
	context.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
}

// POST http://{{host}}/projects/:project-id/submit
func (c *Controller) SubmitProjectSolution(context *gin.Context) {
	var request models.Submission
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.ProjectID = context.Param("project-id")
	// prevent adding when id already present
	for _, submission := range c.Submissions {
		if request.ID == submission.ID || request.ProjectID == submission.ProjectID && request.StudentID == submission.StudentID {
			context.JSON(http.StatusConflict, gin.H{"error": "Submission already present"})
			return
		}
	}
	// update in-memory state
	c.Submissions = append(c.Submissions, request)
	fmt.Println("(SubmitProjectSolution) > In-memory Submission: ", c.Submissions)
	err := kafkaUtils.ProduceMessage(c.Producer, "add", "submission", request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// wait for all messages to be acknowledged
	c.Producer.Flush(15 * 1000)
	context.JSON(http.StatusCreated, gin.H{"message": "Submission created successfully"})

}

// TODO: check whether student and project exists before storing in memory
func (c *Controller) GradeProjectSolution(context *gin.Context) {
	var request models.Grade
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// prevent adding when id or submission_id already present
	for _, grade := range c.Grades {
		if grade.ID == request.ID || grade.SubmissionID == request.SubmissionID {
			context.JSON(http.StatusConflict, gin.H{"error": "Grade with such ID already present or submission already graded"})
			return
		}
	}
	c.Grades = append(c.Grades, request)
	fmt.Println("(GradeProjectSolution) > In-memory Grades: ", c.Grades)
	err := kafkaUtils.ProduceMessage(c.Producer, "add", "grade", request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// wait for all messages to be acknowledged
	c.Producer.Flush(15 * 1000)
	context.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Submission %s graded successfully with %s", request.SubmissionID, request.Grade)})
}
