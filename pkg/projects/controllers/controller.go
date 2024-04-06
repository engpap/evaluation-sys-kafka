package controllers

import (
	kafkaUtils "evaluation-sys-kafka/internal/kafka"
	courseModels "evaluation-sys-kafka/pkg/courses/models"
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
	// In-memory data structures that will be populated through consuming
	Courses []courseModels.Course
}

func (c *Controller) CreateProject(context *gin.Context) {
	var request models.Project
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	courseID := context.Param("course-id")
	request.CourseID = courseID
	// prevent adding when id already present
	for _, project := range c.Projects {
		if project.ID == request.ID {
			context.JSON(http.StatusConflict, gin.H{"error": "Project with such ID already present"})
			return
		}
	}
	// check whether course id exists
	found := false
	for _, course := range c.Courses {
		if course.ID == request.CourseID {
			found = true
			break
		}
	}
	if !found {
		context.JSON(http.StatusBadRequest, gin.H{"error": "You cannot create a project for a course that does not exists"})
		return
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
	// prevent adding when id already present
	request.ProjectID = context.Param("project-id")
	for _, submission := range c.Submissions {
		if request.ID == submission.ID || request.ProjectID == submission.ProjectID && request.StudentID == submission.StudentID {
			context.JSON(http.StatusConflict, gin.H{"error": "Submission already present"})
			return
		}
	}
	courseID := context.Param("course-id")
	// check whether given project is under specified course
	// find project by id
	for _, project := range c.Projects {
		if project.ID == request.ProjectID && courseID != project.CourseID {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Submitting solution for a project whose course id is wrong"})
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
	courseID := context.Param("course-id")
	projectID := context.Param("project-id")
	submissionID := context.Param("submission-id")
	// check that project id and submission id make sense, i.e. the submission id corresponds to a submission of the project id given
	// find the submission by id and check whether its project-id corresponds to the one in the URL
	for _, sub := range c.Submissions {
		if sub.ID == submissionID && sub.ProjectID != projectID {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Grading a submission with a wrong project id"})
			return
		}
	}
	// check if project and course make sense
	// check whether given project is under specified course
	// find project by id
	for _, project := range c.Projects {
		if project.ID == projectID && courseID != project.CourseID {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Grading a solution with a wrong course id"})
			return
		}
	}
	var request models.Grade
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.SubmissionID = submissionID
	// prevent adding when id or submission_id already present
	for _, grade := range c.Grades {
		if grade.ID == request.ID || grade.SubmissionID == submissionID {
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

func (c *Controller) SaveCourseInMemory(data interface{}) {
	if courseMap, ok := data.(map[string]interface{}); ok {
		course := courseModels.Course{
			ID:   fmt.Sprint(courseMap["id"]),
			Name: fmt.Sprint(courseMap["name"]),
		}
		c.Courses = append(c.Courses, course)
		fmt.Println("In-Memory Courses: ", c.Courses)
	} else {
		fmt.Printf("Error: data cannot be converted to Course\n")
	}
}
