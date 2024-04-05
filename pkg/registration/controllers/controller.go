package controllers

import (
	courseModels "evaluation-sys-kafka/pkg/courses/models"
	projectModels "evaluation-sys-kafka/pkg/projects/models"
	"fmt"
)

type Controller struct {
	Courses  []courseModels.Course
	Projects []projectModels.Project
	Grades   []projectModels.Grade
}

func (c *Controller) GetCourses() {
	// TODO: listen on CourseConsumerOutput and updates Courses
}

func (c *Controller) SaveCourseInMemory(data interface{}) {
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
func (c *Controller) SaveProjectInMemory(data interface{}) {
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

func (c *Controller) SaveGradeInMemory(data interface{}) {
	if gradeMap, ok := data.(map[string]interface{}); ok {
		grade := projectModels.Grade{
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
