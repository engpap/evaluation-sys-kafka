package transport

import (
	"eval-sys-registration-service/cmd/internal/controllers"
	"fmt"

	kafkaWrapper "github.com/engpap/kafka-wrapper-go/pkg"

	"github.com/gin-gonic/gin"
)

func Serve() {
	debug := false
	port := "8080"
	if debug {
		port = "8092"
	}
	router := initRouter()
	router.Run(":" + port)
	fmt.Println("Server is running on port " + port)
}

func initRouter() *gin.Engine {
	router := gin.Default()

	producer, err := kafkaWrapper.CreateProducer()
	if err != nil {
		panic(err)
	}

	registrationController := controllers.Controller{Producer: producer}
	kafkaWrapper.SetupCloseProducerHandler(producer)

	// consumes on events created by other services
	go kafkaWrapper.CreateConsumer("course", registrationController.UpdateCourseInMemory, "registration-service")
	go kafkaWrapper.CreateConsumer("project", registrationController.UpdateProjectInMemory, "registration-service")
	go kafkaWrapper.CreateConsumer("grade", registrationController.UpdateGradeInMemory, "registration-service")
	go kafkaWrapper.CreateConsumer("submission", registrationController.UpdateSubmissionInMemory, "registration-service")
	go kafkaWrapper.CreateConsumer("completed", registrationController.UpdateCompletedCoursesInMemory, "registration-service")

	return router
}
