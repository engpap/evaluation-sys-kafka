package common

import (
	"bufio"
	"eval-sys-cli-app/cmd/internal/config"
	"net/http"

	"github.com/fatih/color"
)

func GetCourses() {
	resp, err := http.Get(config.URLs.CourseServiceURL + "/courses")
	if err != nil {
		color.Red("Error fetching courses: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		color.Red("Failed to fetch courses: %s\n", resp.Status)
	} else {
		color.Green("Courses:")
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			color.Green(scanner.Text())
		}
	}
}
