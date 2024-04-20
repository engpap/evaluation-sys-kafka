package controllers

import (
	"eval-sys-project-service/cmd/internal/models"
	"sync"

	kafkaWrapper "github.com/engpap/kafka-wrapper-go/pkg"

	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	mu       sync.Mutex
	Producer *kafka.Producer
	// In-memory data structures that will be populated through consuming
	Courses     []models.Course
	Enrollments []models.Enrollment
	Projects    []models.Project
	Submissions []models.Submission
	Grades      []models.Grade
	Professors  []models.Professor
}

func (c *Controller) CreateProject(context *gin.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
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
	err := kafkaWrapper.ProduceMessage(c.Producer, "add", "project", request)
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
	c.mu.Lock()
	defer c.mu.Unlock()
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
	// allow submission only by enrolled students
	studentEnrolled := false
	for _, enrollment := range c.Enrollments {
		if enrollment.StudentID == request.StudentID && enrollment.CourseID == courseID {
			studentEnrolled = true
			break
		}
	}
	if !studentEnrolled {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Student not enrolled in course"})
		return
	}
	err := kafkaWrapper.ProduceMessage(c.Producer, "add", "submission", request)
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
	c.mu.Lock()
	defer c.mu.Unlock()
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
	// check whether professor exists
	professorFound := false
	for _, professor := range c.Professors {
		if professor.ID == request.ProfessorID {
			professorFound = true
			break
		}
	}
	if !professorFound {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Professor does not exist"})
		return
	}
	err := kafkaWrapper.ProduceMessage(c.Producer, "add", "grade", request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// wait for all messages to be acknowledged
	c.Producer.Flush(15 * 1000)
	context.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Submission %s graded successfully with %s", request.SubmissionID, request.Grade)})
}

//  ------------------------------------------------------------------------------------------
//  ------------------------------------------------------------------------------------------
// CALLBACK FUNCTIONS

func (c *Controller) UpdateCourseInMemory(action_type string, data interface{}) {
	fmt.Printf("UpdateCourseInMemory > action_type: %s\n", action_type)
	c.mu.Lock()
	defer c.mu.Unlock()
	if action_type == "add" {
		c.saveCourseInMemory(data)
	} else if action_type == "delete" {
		c.deleteCourseInMemory(data)
	} else {
		fmt.Printf("Invalid action type: %s\n", action_type)
	}
}

func (c *Controller) saveCourseInMemory(data interface{}) {
	if courseMap, ok := data.(map[string]interface{}); ok {
		course := models.Course{
			ID:   fmt.Sprint(courseMap["id"]),
			Name: fmt.Sprint(courseMap["name"]),
		}
		c.Courses = append(c.Courses, course)
		fmt.Println("In-Memory Courses: ", c.Courses)
	} else {
		fmt.Printf("Error: data cannot be converted to Course\n")
	}
}

func (c *Controller) deleteCourseInMemory(data interface{}) {
	if courseMap, ok := data.(map[string]interface{}); ok {
		courseID := fmt.Sprint(courseMap["id"])
		for i, course := range c.Courses {
			if course.ID == courseID {
				c.Courses = append(c.Courses[:i], c.Courses[i+1:]...)
			}
		}
		fmt.Println("In-Memory Courses: ", c.Courses)
	} else {
		fmt.Printf("Error: data cannot be converted to Course\n")
	}
}

func (c *Controller) UpdateEnrollmentInMemory(action_type string, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if action_type == "add" {
		c.saveEnrollmentInMemory(data)
	} else {
		fmt.Printf("Invalid action type: %s\n", action_type)
	}
}

func (c *Controller) saveEnrollmentInMemory(data interface{}) {
	if enrollmentMap, ok := data.(map[string]interface{}); ok {
		enrollment := models.Enrollment{
			StudentID: fmt.Sprint(enrollmentMap["student_id"]),
			CourseID:  fmt.Sprint(enrollmentMap["course_id"]),
		}
		c.Enrollments = append(c.Enrollments, enrollment)
		fmt.Println("In-Memory Enrollments: ", c.Enrollments)
	} else {
		fmt.Printf("Error: data cannot be converted to Enrollment\n")
	}
}

func (c *Controller) UpdateProjectInMemory(action_type string, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if action_type == "add" {
		c.saveProjectInMemory(data)
	} else {
		fmt.Printf("Invalid action type: %s\n", action_type)
	}
}

func (c *Controller) saveProjectInMemory(data interface{}) {
	if projectMap, ok := data.(map[string]interface{}); ok {
		project := models.Project{
			ID:       fmt.Sprint(projectMap["id"]),
			Name:     fmt.Sprint(projectMap["name"]),
			CourseID: fmt.Sprint(projectMap["course_id"]),
		}
		c.Projects = append(c.Projects, project)
		fmt.Println("In-Memory Projects: ", c.Projects)
	} else {
		fmt.Printf("Error: data cannot be converted to Project\n")
	}
}

func (c *Controller) UpdateSubmissionInMemory(action_type string, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if action_type == "add" {
		c.saveSubmissionInMemory(data)
	} else {
		fmt.Printf("Invalid action type: %s\n", action_type)
	}
}

func (c *Controller) saveSubmissionInMemory(data interface{}) {
	if submissionMap, ok := data.(map[string]interface{}); ok {
		submission := models.Submission{
			ID:        fmt.Sprint(submissionMap["id"]),
			StudentID: fmt.Sprint(submissionMap["student_id"]),
			ProjectID: fmt.Sprint(submissionMap["project_id"]),
			Solution:  fmt.Sprint(submissionMap["solution"]),
		}
		c.Submissions = append(c.Submissions, submission)
		fmt.Println("In-Memory Submissions: ", c.Submissions)
	} else {
		fmt.Printf("Error: data cannot be converted to Submission\n")
	}
}

func (c *Controller) UpdateGradeInMemory(action_type string, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if action_type == "add" {
		c.saveGradeInMemory(data)
	} else {
		fmt.Printf("Invalid action type: %s\n", action_type)
	}
}

func (c *Controller) saveGradeInMemory(data interface{}) {
	if gradeMap, ok := data.(map[string]interface{}); ok {
		grade := models.Grade{
			ID:           fmt.Sprint(gradeMap["id"]),
			SubmissionID: fmt.Sprint(gradeMap["submission_id"]),
			ProfessorID:  fmt.Sprint(gradeMap["professor_id"]),
			Grade:        fmt.Sprint(gradeMap["grade"]),
		}
		c.Grades = append(c.Grades, grade)
		fmt.Println("In-Memory Grades: ", c.Grades)
	} else {
		fmt.Printf("Error: data cannot be converted to Grade\n")
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
		fmt.Println("In-Memory Professors: ", c.Professors)
	} else {
		fmt.Printf("Error: data cannot be converted to Professor\n")
	}
}

// ------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------
// NON-REQUIRED FUNCTIONS BY THE SPECS

func (c *Controller) GetCourseProjects(context *gin.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	courseID := context.Param("course-id")
	var projects []models.Project
	for _, project := range c.Projects {
		if project.CourseID == courseID {
			projects = append(projects, project)
		}
	}
	context.JSON(http.StatusOK, projects)
}

func (c *Controller) GetProjectSubmissions(context *gin.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	courseID := context.Param("course-id")
	projectID := context.Param("project-id")
	// check whether project id is under specified course
	// find project by id
	for _, project := range c.Projects {
		if project.ID == projectID && courseID != project.CourseID {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Getting submissions for a project with a wrong course id"})
			return
		}
	}
	var submissions []models.Submission
	for _, submission := range c.Submissions {
		if submission.ProjectID == projectID {
			submissions = append(submissions, submission)
		}
	}
	context.JSON(http.StatusOK, submissions)
}

func (c *Controller) GetSubmissionGrades(context *gin.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	courseID := context.Param("course-id")
	projectID := context.Param("project-id")
	submissionID := context.Param("submission-id")
	// check whether project id and submission id make sense, i.e. the submission id corresponds to a submission of the project id given
	// find the submission by id and check whether its project-id corresponds to the one in the URL
	for _, sub := range c.Submissions {
		if sub.ID == submissionID && sub.ProjectID != projectID {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Getting grades for a submission with a wrong project id"})
			return
		}
	}
	// check if project and course make sense
	// check whether given project is under specified course
	// find project by id
	for _, project := range c.Projects {
		if project.ID == projectID && courseID != project.CourseID {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Getting grades for a solution with a wrong course id"})
			return
		}
	}
	var grades []models.Grade
	for _, grade := range c.Grades {
		if grade.SubmissionID == submissionID {
			grades = append(grades, grade)
		}
	}
	context.JSON(http.StatusOK, grades)
}
