package admin

import (
	"evaluation-sys-kafka/cmd/cli-frontend/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

// create-course --id=<id> --name=<name>
// args: ["--id=<id>", "--name=<name>"]
func CreateCourse(args []string) {
	var id, name string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--id=") {
			id = strings.TrimPrefix(arg, "--id=")
		} else if strings.HasPrefix(arg, "--name=") {
			name = strings.TrimPrefix(arg, "--name=")
		}
	}
	if id == "" || name == "" {
		color.Yellow("Usage: create-course --id=<id> --name=<name>\n")
		return
	}
	body := strings.NewReader(fmt.Sprintf(`{"id": "%s", "name": "%s"}`, id, name))
	resp, err := http.Post(config.CourseServiceURL+"/courses/create", "application/json", body)
	if err != nil {
		color.Red("Error creating course: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		color.Red("Failed to create course: %s\n", resp.Status)
		return
	} else {
		color.Green("Course created successfully\n")
	}
}

// delete-course --id=<id>
func DeleteCourse(args []string) {
	var courseID string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--id=") {
			courseID = strings.TrimPrefix(arg, "--id=")
		}
	}
	if courseID == "" {
		color.Yellow("Usage: delete-course --id=<id>\n")
		return
	}
	// endpoint is -> /courses/:course-id/delete
	// Creating the request
	req, err := http.NewRequest(http.MethodDelete, config.CourseServiceURL+"/courses/"+courseID+"/delete", nil)
	if err != nil {
		color.Red("Error creating request: %v\n", err)
		return
	}

	// Sending the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		color.Red("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		color.Red("Failed to delete course: %s\n", resp.Status)
		return
	} else {
		color.Green("Course %s deleted successfully\n", courseID)
	}
}

// create-student --id=<id>
func CreateStudent(args []string) {
	var id string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--id=") {
			id = strings.TrimPrefix(arg, "--id=")
		}
	}
	if id == "" {
		color.Yellow("Usage: create-student --id=<id>\n")
		return
	}
	body := strings.NewReader(fmt.Sprintf(`{"id": "%s"}`, id))
	resp, err := http.Post(config.UserServiceURL+"/users/student/create", "application/json", body)
	if err != nil {
		color.Red("Error creating student: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		color.Red("Failed to create student: %s\n", resp.Status)
		return
	} else {
		color.Green("Student created successfully\n")
	}

}

// create-professor --id=<id>
func CreateProfessor(args []string) {
	var id string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--id=") {
			id = strings.TrimPrefix(arg, "--id=")
		}
	}
	if id == "" {
		color.Yellow("Usage: create-professor --id=<id>\n")
		return
	}
	body := strings.NewReader(fmt.Sprintf(`{"id": "%s"}`, id))
	resp, err := http.Post(config.UserServiceURL+"/users/professor/create", "application/json", body)
	if err != nil {
		color.Red("Error creating professor: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		color.Red("Failed to create professor: %s\n", resp.Status)
		return
	} else {
		color.Green("Professor created successfully\n")
	}
}
