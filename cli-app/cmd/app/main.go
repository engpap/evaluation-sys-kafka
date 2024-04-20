package main

import (
	"bufio"
	admin "eval-sys-cli-app/cmd/internal/controllers/admin"
	common "eval-sys-cli-app/cmd/internal/controllers/common"
	professor "eval-sys-cli-app/cmd/internal/controllers/professor"
	student "eval-sys-cli-app/cmd/internal/controllers/student"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fatih/color"
)

var userID string

// Run: go run cmd/cli-frontend/main.go
func main() {
	interrupt := make(chan os.Signal, 1) // channel to listen for interrupt signal (Ctrl+C)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	reader := bufio.NewReader(os.Stdin) // buffered reader for standard input

	var handleInput func(string, []string)
	var role string
	var err error

	color.Yellow("Select role (admin, student, professor): ")
	for {
		role, err = reader.ReadString('\n')
		if err != nil {
			color.Red("Error reading role: %v\n", err)
			return
		}
		role = strings.TrimSpace(role)
		if role == "admin" {
			handleInput = handleAdminInput
			break
		} else if role == "student" {
			handleInput = handleStudentInput
			break
		} else if role == "professor" {
			handleInput = handleProfessorInput
			break
		} else {
			color.Red("Invalid role. Try again")
		}
		// check for interrupt signal in a non-blocking way
		select {
		case <-interrupt:
			color.Magenta("Interrupt signal received. Exiting...")
			return
		default: // if no interrupt signal was received, continue as normal
		}
	}
	// ask the user his credentials (userID)
	color.Yellow("Enter your %s id: ", role)
	userID, err = reader.ReadString('\n')
	if err != nil {
		color.Red("Error reading %s id: %v\n", role, err)
		return
	}
	userID = strings.TrimSpace(userID)
	select {
	case <-interrupt:
		color.Magenta("Interrupt signal received. Exiting...")
		return
	default: // if no interrupt signal was received, continue as normal
	}
	color.Cyan("Logged in as %s with ID %s", role, userID)
	fmt.Println()
	// print possible commands
	color.Cyan("Here's a list of commands you can use:")
	if role == "admin" {
		color.Cyan("create-course --id=<id> --name=<name>")
		color.Cyan("get-courses")
		color.Cyan("delete-course --id=<id>")
		color.Cyan("create-student --id=<id>")
		color.Cyan("create-professor --id=<id>")
	} else if role == "student" {
		color.Cyan("get-courses")
		color.Cyan("enroll --course-id=<id>")
		color.Cyan("submit-solution --course-id=<id> --project-id=<id> --submission-id=<id> --solution=<solution>")
		color.Cyan("get-course-projects --course-id=<id>")
		color.Cyan("get-project-submissions --course-id=<id> --project-id=<id>")
		color.Cyan("get-submission-grades --course-id=<id> --project-id=<id> --submission-id=<id>")
	} else if role == "professor" {
		color.Cyan("get-courses")
		color.Cyan("create-project --id=<id> --course-id=<course-id> --name=<project-name>")
		color.Cyan("get-sub --course-id=<id> --project-id=<id>")
		color.Cyan("grade --course-id=<id> --proj-id=<id> --sub-id=<id> --grade-id=<grade-id> --grade=<grade>")
	} else {
		color.Red("Invalid role. Try again")
	}
	fmt.Println()
	// loop indefinitely until interrupted with Ctrl+C
	for {
		color.New(color.FgGreen).Fprint(os.Stdout, "> ") // otherwise there's a newline after the prompt
		input, err := reader.ReadString('\n')            // read input until newline
		if err != nil {
			color.Red("Error reading input: %v\n", err)
			continue
		}
		input = strings.TrimSpace(input)
		// if the input is "exit" or "quit", break the loop
		if input == "exit" || input == "quit" {
			color.Magenta("Exiting...")
			break
		}
		// split the input into command and arguments
		parts := strings.Fields(input)
		if len(parts) == 0 {
			return
		}
		command := parts[0]
		args := parts[1:]
		handleInput(command, args)
		// check for interrupt signal in a non-blocking way
		select {
		case <-interrupt:
			color.Magenta("Interrupt signal received. Exiting...")
			return
		default: // if no interrupt signal was received, continue as normal
		}
	}
}

func handleAdminInput(command string, args []string) {
	switch command {
	case "create-course":
		admin.CreateCourse(args)
	case "get-courses":
		common.GetCourses()
	case "delete-course":
		admin.DeleteCourse(args)
	case "create-student":
		admin.CreateStudent(args)
	case "create-professor":
		admin.CreateProfessor(args)
	default:
		color.Red("Unknown command: %s\n", command)
	}
}

func handleStudentInput(command string, args []string) {
	switch command {
	case "get-courses":
		common.GetCourses()
	case "enroll":
		student.EnrollStudentInCourse(userID, args)
	case "submit-solution":
		student.SubmitProjectSolution(userID, args)
	case "get-course-projects":
		student.GetCourseProjects(args)
	case "get-project-submissions":
		student.GetProjectSubmissions(args)
	case "get-submission-grades":
		student.GetSubmissionGrades(args)
	default:
		color.Red("Unknown command: %s\n", command)
	}
}

func handleProfessorInput(command string, args []string) {
	switch command {
	case "get-courses":
		common.GetCourses()
	case "create-project":
		professor.CreateProject(args)
	case "get-sub":
		professor.GetProjectSubmissions(args)
	case "grade":
		professor.GradeProjectSolution(userID, args)
	default:
		color.Red("Unknown command: %s\n", command)
	}
}
