package student

import (
	"evaluation-sys-kafka/cmd/cli-frontend/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

// expected json body by backend -> {"student_id": "id", "course_id": "id"}
func EnrollStudentInCourse(studentID string, args []string) {
	var courseID string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--course-id=") {
			courseID = strings.TrimPrefix(arg, "--course-id=")
		}
	}
	if courseID == "" {
		color.Yellow("Usage: enroll --course-id=<id>\n")
		return
	}
	body := strings.NewReader(fmt.Sprintf(`{"student_id": "%s", "course_id": "%s"}`, studentID, courseID))
	resp, err := http.Post(config.CourseServiceURL+"/courses/enroll", "application/json", body)
	if err != nil {
		color.Red("Error enrolling student in course: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		color.Red("Failed to enroll student in course: %s\n", resp.Status)
		return
	} else {
		color.Green("Student enrolled in course successfully\n")
	}

}
