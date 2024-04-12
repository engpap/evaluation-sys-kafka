package transport

import (
	kafkaUtils "evaluation-sys-kafka/internal/kafka"
	"evaluation-sys-kafka/pkg/courses/controllers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Serve() {
	port := "8080"
	router := initRouter()
	router.Run(":" + port)
	fmt.Println("Server is running on port " + port)
}

func initRouter() *gin.Engine {
	router := gin.Default()

	producer, err := kafkaUtils.CreateProducer()
	if err != nil {
		panic(err)
	}
	// producer setup
	courseController := controllers.Controller{Producer: producer}
	kafkaUtils.SetupCloseProducerHandler(producer)
	// consumers setup
	go kafkaUtils.CreateConsumer("student", courseController.UpdateStudentInMemory)

	// routes
	router.GET("/courses", courseController.GetCourses)                        // FE done 4 everybody
	router.POST("/courses/create", courseController.CreateCourse)              // FE done for admin
	router.DELETE("/courses/:course-id/delete", courseController.DeleteCourse) // FE done for admin
	router.POST("/courses/enroll", courseController.EnrollStudentInCourse)     // FE done for stud

	return router
}
