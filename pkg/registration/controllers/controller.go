package controllers

import "evaluation-sys-kafka/pkg/courses/models"

type Controller struct {
	CourseConsumerOutput []interface{}
	Courses              []models.Course
}

func (c *Controller) GetCourses() {
	// TODO: listen on CourseConsumerOutput and updates Courses
}
