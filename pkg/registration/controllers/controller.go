package controllers

import (
	coursesModels "evaluation-sys-kafka/pkg/courses/models"
	projectsModels "evaluation-sys-kafka/pkg/projects/models"
)

type Controller struct {
	CourseConsumerOutput  []interface{}
	Courses               []coursesModels.Course
	ProjectConsumerOutput []interface{}
	Projects              []projectsModels.Project
	GradeConsumerOutput   []interface{}
	Grades                []projectsModels.Grade
}

func (c *Controller) GetCourses() {
	// TODO: listen on CourseConsumerOutput and updates Courses
}
