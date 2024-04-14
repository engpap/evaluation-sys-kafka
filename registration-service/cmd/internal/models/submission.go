package models

type Submission struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	StudentID string `json:"student_id"`
	Solution  string `json:"solution"`
}
