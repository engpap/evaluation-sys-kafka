package controllers

import (
	courseModels "evaluation-sys-kafka/pkg/courses/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	ConsumerOutput []interface{}
}

func (c *Controller) GetCourses(context *gin.Context) {
	var courses []courseModels.Course // Initialize a slice to hold Course instances

	for index, element := range c.ConsumerOutput {
		//fmt.Printf("Index: %d, Value: %v, Type: %T\n", index, element, element)
		if courseMap, ok := element.(map[string]interface{}); ok {
			course := courseModels.Course{
				ID:   fmt.Sprint(courseMap["id"]),
				Name: fmt.Sprint(courseMap["name"]),
			}
			courses = append(courses, course)
		} else {
			fmt.Printf("Error: element at index %d cannot be converted to Course\n", index)
		}
	}
	context.JSON(http.StatusOK, gin.H{"courses": courses})
}

/*

func (c *Controller) GetCourses(context *gin.Context) {
	fmt.Printf(">>> " + string(c.ConsumerOutput))

	// Try unmarshalling into an array of Courses first
	var courses []courseModels.Course
	err := json.Unmarshal(c.ConsumerOutput, &courses)

	// If unmarshalling into an array fails, try unmarshalling into a single Course object
	if err != nil {
		var singleCourse courseModels.Course
		err = json.Unmarshal(c.ConsumerOutput, &singleCourse)
		// If this succeeds, append the single course to the courses slice
		if err == nil {
			courses = append(courses, singleCourse)
		}
	}

	// If both attempts fail, return an error
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("GetCourses\n")
	fmt.Printf("ConsumerOutput: %s\n", c.ConsumerOutput)

	context.JSON(http.StatusOK, gin.H{"courses": courses})
}*/
