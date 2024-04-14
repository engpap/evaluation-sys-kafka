package models

type Project struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	CourseID string `json:"course_id"`
}
