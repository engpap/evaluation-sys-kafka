package models

type Submission struct {
	ID        string `json:"id"`
	ProjectID string `json:"-"` // "-" means this field is ignored by the JSON parser
	StudentID string `json:"student_id"`
	Solution  string `json:"solution"`
}
