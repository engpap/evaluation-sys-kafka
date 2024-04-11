package controllers

import (
	"errors"
	courseModels "evaluation-sys-kafka/pkg/courses/models"
	projectModels "evaluation-sys-kafka/pkg/projects/models"
	"fmt"
	"os"
)

type Controller struct {
	Courses     []courseModels.Course
	Projects    []projectModels.Project
	Grades      []projectModels.Grade
	Submissions []projectModels.Submission
}

func (c *Controller) UpdateCourseInMemory(action_type string, data interface{}) {
	if action_type == "add" {
		c.saveCourseInMemory(data)
	} else {
		fmt.Printf("Error: action type %s not supported\n", action_type)
	}
}

func (c *Controller) saveCourseInMemory(data interface{}) {
	if courseMap, ok := data.(map[string]interface{}); ok {
		course := courseModels.Course{
			ID:   fmt.Sprint(courseMap["id"]),
			Name: fmt.Sprint(courseMap["name"]),
		}
		c.Courses = append(c.Courses, course)
		fmt.Println("In-Memory Course: ", c.Courses)
	} else {
		fmt.Printf("Error: data cannot be converted to Course\n")
	}
}

func (c *Controller) UpdateProjectInMemory(action_type string, data interface{}) {
	if action_type == "add" {
		c.saveProjectInMemory(data)
	} else {
		fmt.Printf("Error: action type %s not supported\n", action_type)
	}
}

func (c *Controller) saveProjectInMemory(data interface{}) {
	if projectMap, ok := data.(map[string]interface{}); ok {
		project := projectModels.Project{
			ID:       fmt.Sprint(projectMap["id"]),
			Name:     fmt.Sprint(projectMap["name"]),
			CourseID: fmt.Sprint(projectMap["course_id"]),
		}
		c.Projects = append(c.Projects, project)
		fmt.Println("In-Memory Project: ", c.Projects)
	} else {
		fmt.Printf("Error: data cannot be converted to Project\n")
	}
}

func (c *Controller) UpdateSubmissionInMemory(action_type string, data interface{}) {
	if action_type == "add" {
		c.saveSubmissionInMemory(data)
	} else {
		fmt.Printf("Error: action type %s not supported\n", action_type)
	}
}

func (c *Controller) saveSubmissionInMemory(data interface{}) {
	if submissionMap, ok := data.(map[string]interface{}); ok {
		submission := projectModels.Submission{
			ID:        fmt.Sprint(submissionMap["id"]),
			ProjectID: fmt.Sprint(submissionMap["project_id"]),
			StudentID: fmt.Sprint(submissionMap["student_id"]),
			Solution:  fmt.Sprint(submissionMap["solution"]),
		}
		c.Submissions = append(c.Submissions, submission)
		fmt.Println("In-Memory Submissions: ", c.Submissions)
	} else {
		fmt.Printf("Error: data cannot be converted to Submission\n")
	}
}

func (c *Controller) UpdateGradeInMemory(action_type string, data interface{}) {
	if action_type == "add" {
		c.saveGradeInMemory(data)
	} else {
		fmt.Printf("Error: action type %s not supported\n", action_type)
	}
}

func (c *Controller) saveGradeInMemory(data interface{}) {
	if gradeMap, ok := data.(map[string]interface{}); ok {
		grade := projectModels.Grade{
			ID:           fmt.Sprint(gradeMap["id"]),
			SubmissionID: fmt.Sprint(gradeMap["submission_id"]),
			ProfessorID:  fmt.Sprint(gradeMap["professor_id"]),
			Grade:        fmt.Sprint(gradeMap["grade"]),
		}
		c.Grades = append(c.Grades, grade)
		fmt.Println("In-Memory Grades: ", c.Grades)
		// check if the course is completed for the student
		// get student id from the submission
		studentID, err := c.getStudentIDofSubmission(grade.SubmissionID)
		if err != nil {
			fmt.Println("Failed to check course completion: ", err.Error())
			os.Exit(1)
		}
		// get the course id from the submission
		courseID, err := c.getCourseIDofSubmission(grade.SubmissionID)
		if err != nil {
			fmt.Println("Failed to check course completion: ", err.Error())
			os.Exit(1)
		}
		if c.isCourseCompleted(studentID, courseID) {
			fmt.Println("Course is completed for student: ", studentID)
		}
	} else {
		fmt.Printf("Error: data cannot be converted to Grade\n")
	}
}

// A course is completed for a student if the student delivered all projects
// for that course and the sum of the grades is suffcient.
// ASSUMPTION #1: the sum of the grades is sufficient when sum > 50
// ASSUMPTION #2: each course have several projects; the max grade for a course is 100; each project's grade of a course is cumulated.
func (c *Controller) isCourseCompleted(studentID string, courseID string) bool {
	// get all the project of the specified course (courseID)
	var courseProjects []projectModels.Project
	for _, project := range c.Projects {
		if project.CourseID == courseID {
			courseProjects = append(courseProjects, project)
		}
	}
	// OBSERVATION: golang does not have sets, the way to do it is with a map of string -> bool
	// create set of project ids of the specified course (courseID) --- for efficient look-up
	projectIDsOfCourse := make(map[string]bool)
	for _, project := range courseProjects {
		projectIDsOfCourse[project.ID] = true
	}
	// get the graded submissions of the specified student (studentID)
	var studentGrades []projectModels.Grade
	// create set of project ids of the graded submissions of the specified student (studentID)
	studentProjectIDsOfGradedSubmissions := make(map[string]bool)
	for _, grade := range c.Grades {
		if c.isGradedSubmissionOfStudent(grade.SubmissionID, studentID) {
			studentGrades = append(studentGrades, grade)
			// store the project id of the graded submission (for later computation)
			projectIDofGradedSubmission, err := c.getProjectIDofSubmission(grade.SubmissionID)
			if err != nil {
				// exit from program when error is returned
				fmt.Println("Failed to check course completion: ", err.Error())
				os.Exit(1)
			}
			studentProjectIDsOfGradedSubmissions[projectIDofGradedSubmission] = true
		}
	}
	// if there's a project among the the projects of courseID
	// that is not inside the graded project's submissions => course is not completed, return false
	for _, project := range courseProjects {
		// if there not exists a course project in the projects of the student ==> not completed course
		_, exists := studentProjectIDsOfGradedSubmissions[project.ID]
		if !exists {
			return false
		}
	}
	// compute sum of grades & check whether is > 50
	// consider only the grades of the student's submissions of the specified course
	// so the grades whose submission id is in the projectIDsOfCourse
	var studentSubmissionGradesOfProjectOfSpecifiedCourse []projectModels.Grade
	for _, grade := range studentGrades {
		projectIDofGradedSubmission, err := c.getProjectIDofSubmission(grade.SubmissionID)
		if err != nil {
			// exit from program when error is returned
			fmt.Println("Failed to check course completion: ", err.Error())
			os.Exit(1)
		}
		// check if the project id of the graded submission is in the projectIDsOfCourse
		_, exists := projectIDsOfCourse[projectIDofGradedSubmission]
		if exists {
			studentSubmissionGradesOfProjectOfSpecifiedCourse = append(studentSubmissionGradesOfProjectOfSpecifiedCourse, grade)
		}
	}
	// compute sum of grades
	sumOfGrades := 0
	for _, grade := range studentSubmissionGradesOfProjectOfSpecifiedCourse {
		gradeValue := 0
		fmt.Sscanf(grade.Grade, "%d", &gradeValue)
		sumOfGrades += gradeValue
	}
	fmt.Printf("(isCourseCompleted) > Running sum of grades of student %s for course %s is: %d\n", studentID, courseID, sumOfGrades)
	// check if the sum of grades is sufficient
	if sumOfGrades > 50 {
		return true
	} else {
		return false
	}
}

func (c *Controller) isGradedSubmissionOfStudent(submissionID string, studentID string) bool {
	// find submission by id and check if student id is equal to the given id
	for _, submission := range c.Submissions {
		if submission.ID == submissionID && submission.StudentID == studentID {
			return true
		}
	}
	return false
}

func (c *Controller) getProjectIDofSubmission(submissionID string) (string, error) {
	// find submission by id and return the project id stored inside it
	for _, submission := range c.Submissions {
		if submission.ID == submissionID {
			return submission.ProjectID, nil
		}
	}
	return "", errors.New("could not find project id given the submission id")
}

func (c *Controller) getStudentIDofSubmission(submissionID string) (string, error) {
	// find submission by id and return the student id stored inside it
	for _, submission := range c.Submissions {
		if submission.ID == submissionID {
			return submission.StudentID, nil
		}
	}
	return "", errors.New("could not find student id given the submission id")
}

func (c *Controller) getCourseIDofSubmission(submissionID string) (string, error) {
	// get project id of submission using submission id
	projectID, err := c.getProjectIDofSubmission(submissionID)
	if err != nil {
		return "", err
	}
	// get course id of a project using project id
	for _, project := range c.Projects {
		if project.ID == projectID {
			return project.CourseID, nil
		}
	}
	return "", errors.New("could not course id given the submission id")
}
