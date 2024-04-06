package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	userServiceURL    = "http://localhost:8083"
	projectServiceURL = "http://localhost:8084"
	courseServiceURL  = "http://localhost:8080"
)

func main() {
	// channel to listen for interrupt signal (Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	// buffered reader for standard input
	reader := bufio.NewReader(os.Stdin)
	// loop indefinitely until interrupted with Ctrl+C
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n') // read input until newline
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}
		// trim the input
		// if the input is "exit" or "quit", break the loop
		input = strings.TrimSpace(input)
		if input == "exit" || input == "quit" {
			fmt.Println("Exiting...")
			break
		}
		handleInput(input)
		// Check for interrupt signal in a non-blocking way
		select {
		case <-interrupt:
			fmt.Println("\nReceived an interrupt, stopping...")
			return
		default:
			// If no interrupt signal was received, continue as normal
		}
	}
}

func handleInput(input string) {
	// Split the input into command and arguments
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "create-course":
		createCourse(args)
	default:
		fmt.Println("Unknown command:", command)
	}
}

// create-course --id=<id> --name=<name>
// args: ["--id=<id>", "--name=<name>"]
func createCourse(args []string) {
	var id, name string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--id=") {
			id = strings.TrimPrefix(arg, "--id=")
		} else if strings.HasPrefix(arg, "--name=") {
			name = strings.TrimPrefix(arg, "--name=")
		}
	}
	if id == "" || name == "" {
		fmt.Println("Usage: create-course --id=<id> --name=<name>")
		return
	}
	body := strings.NewReader(fmt.Sprintf(`{"id": "%s", "name": "%s"}`, id, name))
	resp, err := http.Post(courseServiceURL+"/courses/create", "application/json", body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating course: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		fmt.Fprintf(os.Stderr, "Failed to create course: %s\n", resp.Status)
		return
	} else {
		fmt.Println("Course created successfully")
	}
}
