package models

type Grade struct {
	ID           string `json:"id"`
	SubmissionID string `json:"submission_id"`
	ProfessorID  string `json:"professor_id"`
	Grade        string `json:"grade"`
}
