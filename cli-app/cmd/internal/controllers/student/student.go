package student

import (
	"bufio"
	"eval-sys-cli-app/cmd/internal/config"
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
	body := strings.NewReader(fmt.Sprintf(`{"student_id": "%s"}`, studentID))
	resp, err := http.Post(config.URLs.CourseServiceURL+"/courses/"+courseID+"/enroll", "application/json", body)
	if err != nil {
		color.Red("Error enrolling student in course: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		color.Red("Failed to enroll student in course: %s\n", resp.Status)
	} else {
		color.Green("Student enrolled in course successfully\n")
	}

}

// POST("/courses/:course-id/projects/:project-id/submit")
// JSON body: {"id": "submission-id", "student_id": "student-id", "solution": "solution"}
func SubmitProjectSolution(studentID string, args []string) {
	var courseID, projectID, submissionID, solution string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--course-id=") {
			courseID = strings.TrimPrefix(arg, "--course-id=")
		}
		if strings.HasPrefix(arg, "--project-id=") {
			projectID = strings.TrimPrefix(arg, "--project-id=")
		}
		if strings.HasPrefix(arg, "--submission-id=") {
			submissionID = strings.TrimPrefix(arg, "--submission-id=")
		}
		if strings.HasPrefix(arg, "--solution=") {
			solution = strings.TrimPrefix(arg, "--solution=")
		}
	}
	if courseID == "" || projectID == "" || solution == "" {
		color.Yellow("Usage: submit --course-id=<id> --project-id=<id> --submission-id=<id> --solution=<solution>\n")
		return
	}
	body := strings.NewReader(fmt.Sprintf(`{"id": "%s", "student_id": "%s", "solution": "%s"}`, submissionID, studentID, solution))
	resp, err := http.Post(config.URLs.ProjectServiceURL+"/courses/"+courseID+"/projects/"+projectID+"/submit", "application/json", body)
	if err != nil {
		color.Red("Error submitting project solution: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		color.Red("Failed to submit project solution: %s\n", resp.Status)
	} else {
		color.Green("Project solution submitted successfully with ID %s\n", submissionID)
	}
}

func GetCourseProjects(args []string) {
	var courseID string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--course-id=") {
			courseID = strings.TrimPrefix(arg, "--course-id=")
		}
	}
	if courseID == "" {
		color.Yellow("Usage: get-projects --course-id=<id>\n")
		return
	}
	resp, err := http.Get(config.URLs.ProjectServiceURL + "/courses/" + courseID + "/projects")
	if err != nil {
		color.Red("Error getting projects: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		color.Red("Failed to get projects: %s\n", resp.Status)
		return
	} else {
		color.Green("Projects of the course with ID %s:", courseID)
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			color.Green(scanner.Text())
		}
	}
}

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
		color.Yellow("Usage: get-submissions --course-id=<id> --project-id=<id>\n")
		return
	}
	resp, err := http.Get(config.URLs.ProjectServiceURL + "/courses/" + courseID + "/projects/" + projectID + "/submissions")
	if err != nil {
		color.Red("Error getting submissions: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		color.Red("Failed to get submissions: %s\n", resp.Status)
		return
	} else {
		color.Green("Submissions for the project ID %s and course ID %s:", projectID, courseID)
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			color.Green(scanner.Text())
		}
	}
}

func GetSubmissionGrades(args []string) {
	var courseID, projectID, submissionID string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--course-id=") {
			courseID = strings.TrimPrefix(arg, "--course-id=")
		}
		if strings.HasPrefix(arg, "--project-id=") {
			projectID = strings.TrimPrefix(arg, "--project-id=")
		}
		if strings.HasPrefix(arg, "--submission-id=") {
			submissionID = strings.TrimPrefix(arg, "--submission-id=")
		}
	}
	if courseID == "" || projectID == "" || submissionID == "" {
		color.Yellow("Usage: get-grades --course-id=<id> --project-id=<id> --submission-id=<id>\n")
		return
	}
	resp, err := http.Get(config.URLs.ProjectServiceURL + "/courses/" + courseID + "/projects/" + projectID + "/submissions/" + submissionID + "/grades")
	if err != nil {
		color.Red("Error getting grades: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		color.Red("Failed to get grades: %s\n", resp.Status)
		return
	} else {
		color.Green("Grade of submission ID %s of the project ID %s and course ID %s:", submissionID, projectID, courseID)
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			color.Green(scanner.Text())
		}
	}
}

func GetCompletedCourses() {
	resp, err := http.Get(config.URLs.UserServiceURL + "/users/student/completed-courses")
	if err != nil {
		color.Red("Error getting completed courses: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		color.Red("Failed to get completed courses: %s\n", resp.Status)
		return
	} else {
		color.Green("Completed courses:")
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			color.Green(scanner.Text())
		}
	}
}
