package professor

import (
	"bufio"
	"eval-sys-cli-app/cmd/internal/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

// POST("/courses/:course-id/projects/create")
// JSON body: {"id": "project-id", "name": "project-name"}
func CreateProject(args []string) {
	var courseID, projectID, projectName string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--id=") {
			projectID = strings.TrimPrefix(arg, "--id=")
		}
		if strings.HasPrefix(arg, "--course-id=") {
			courseID = strings.TrimPrefix(arg, "--course-id=")
		}
		if strings.HasPrefix(arg, "--name=") {
			projectName = strings.TrimPrefix(arg, "--name=")
		}
	}
	if courseID == "" || projectID == "" || projectName == "" {
		color.Yellow("Usage: create-project --id=<id> --course-id=<course-id> --name=<project-name>\n")
		return
	}
	body := strings.NewReader(fmt.Sprintf(`{"id": "%s", "name": "%s"}`, projectID, projectName))
	resp, err := http.Post(config.URLs.ProjectServiceURL+"/courses/"+courseID+"/projects/create", "application/json", body)
	if err != nil {
		color.Red("Error creating project: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		color.Red("Failed to create project: %s\n", resp.Status)
		return
	} else {
		color.Green("Project created successfully\n")
	}
}

// POST("/courses/:course-id/projects/:project-id/submissions/:submission-id/grade")
// JSON body: {"id" : "grade-id", "professor_id" : "prof-id", "grade": "grade"}
func GradeProjectSolution(professorID string, args []string) {
	var courseID, projectID, submissionID, gradeID, grade string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--course-id=") {
			courseID = strings.TrimPrefix(arg, "--course-id=")
		}
		if strings.HasPrefix(arg, "--proj-id=") {
			projectID = strings.TrimPrefix(arg, "--proj-id=")
		}
		if strings.HasPrefix(arg, "--sub-id=") {
			submissionID = strings.TrimPrefix(arg, "--sub-id=")
		}
		if strings.HasPrefix(arg, "--grade-id=") {
			gradeID = strings.TrimPrefix(arg, "--grade-id=")
		}
		if strings.HasPrefix(arg, "--grade=") {
			grade = strings.TrimPrefix(arg, "--grade=")
		}
	}
	if courseID == "" || projectID == "" || submissionID == "" || gradeID == "" || grade == "" {
		color.Yellow("Usage: grade --course-id=<id> --proj-id=<id> --sub-id=<id> --id=<grade-id> --grade=<grade>\n")
		return
	}
	body := strings.NewReader(fmt.Sprintf(`{"id": "%s", "professor_id": "%s", "grade": "%s"}`, gradeID, professorID, grade))
	resp, err := http.Post(config.URLs.ProjectServiceURL+"/courses/"+courseID+"/projects/"+projectID+"/submissions/"+submissionID+"/grade", "application/json", body)
	if err != nil {
		color.Red("Error grading project: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		color.Red("Failed to grade project: %s\n", resp.Status)
		return
	} else {
		color.Green("Project graded successfully\n")
	}
}

// get-subs
// /courses/:course-id/projects/:project-id/submissions
func GetProjectSubmissions(args []string) {
	var courseID, projectID string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--course-id=") {
			courseID = strings.TrimPrefix(arg, "--course-id=")
		}
		if strings.HasPrefix(arg, "--project-id=") {
			projectID = strings.TrimPrefix(arg, "--project-id=")
		}
	}
	if courseID == "" || projectID == "" {
		color.Yellow("Usage: get-sub --course-id=<id> --project-id=<id>\n")
		return
	}
	resp, err := http.Get(config.URLs.ProjectServiceURL + "/courses/" + courseID + "/projects/" + projectID + "/submissions")
	if err != nil {
		color.Red("Error fetching submissions: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		color.Red("Failed to fetch submissions: %s\n", resp.Status)
	} else {
		color.Green("Submissions from students:")
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			color.Green(scanner.Text())
		}
	}
}
