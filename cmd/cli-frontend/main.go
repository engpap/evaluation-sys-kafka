package main

import (
	"bufio"
	"evaluation-sys-kafka/cmd/cli-frontend/controllers/admin"
	"evaluation-sys-kafka/cmd/cli-frontend/controllers/common"
	"evaluation-sys-kafka/cmd/cli-frontend/controllers/student"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fatih/color"
)

var UserID string

// Run: go run cmd/cli-frontend/main.go
func main() {
	// channel to listen for interrupt signal (Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	// buffered reader for standard input
	reader := bufio.NewReader(os.Stdin)

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
			color.Cyan("Here's a list of commands you can use:")
			color.Cyan("create-course --id=<id> --name=<name>")
			color.Cyan("get-courses")
			color.Cyan("delete-course --id=<id>")
			color.Cyan("create-student --id=<id> --name=<name>")
			color.Cyan("create-professor --id=<id> --name=<name>")
			handleInput = handleAdminInput
			break
		} else if role == "student" {
			color.Cyan("Here's a list of commands you can use:")
			color.Cyan("get-courses")
			color.Cyan("enroll --course-id=<id>")
			handleInput = handleStudentInput
			break
		} else if role == "professor" {
			color.Cyan("Here's a list of commands you can use:")
			color.Cyan("get-courses")
			handleInput = handleProfessorInput
			break
		} else {
			color.Red("Invalid role. Try again")
		}
		// Check for interrupt signal in a non-blocking way
		select {
		case <-interrupt:
			color.Magenta("Interrupt signal received. Exiting...")
			return
		default:
			// If no interrupt signal was received, continue as normal
		}
	}
	// Ask the user his credentials (userID)
	color.Yellow("Enter your %s id: ", role)
	UserID, err = reader.ReadString('\n')
	if err != nil {
		color.Red("Error reading %s id: %v\n", role, err)
		return
	}
	UserID = strings.TrimSpace(UserID)
	select {
	case <-interrupt:
		color.Magenta("Interrupt signal received. Exiting...")
		return
	default:
		// If no interrupt signal was received, continue as normal
	}

	color.Cyan("Logged in as %s with ID %s", role, UserID)
	// loop indefinitely until interrupted with Ctrl+C
	for {
		color.New(color.FgGreen).Fprint(os.Stdout, "> ") // otherwise there's a newline after the prompt
		input, err := reader.ReadString('\n')            // read input until newline
		if err != nil {
			color.Red("Error reading input: %v\n", err)
			continue
		}
		// trim the input
		// if the input is "exit" or "quit", break the loop
		input = strings.TrimSpace(input)
		if input == "exit" || input == "quit" {
			color.Magenta("Exiting...")
			break
		}
		// Split the input into command and arguments
		parts := strings.Fields(input)
		if len(parts) == 0 {
			return
		}
		command := parts[0]
		args := parts[1:]
		handleInput(command, args)
		// Check for interrupt signal in a non-blocking way
		select {
		case <-interrupt:
			color.Magenta("Interrupt signal received. Exiting...")
			return
		default:
			// If no interrupt signal was received, continue as normal
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
		student.EnrollStudentInCourse(UserID, args)
	default:
		color.Red("Unknown command: %s\n", command)
	}
}

func handleProfessorInput(command string, args []string) {
	switch command {
	case "get-courses":
		common.GetCourses()
	default:
		color.Red("Unknown command: %s\n", command)
	}
}
