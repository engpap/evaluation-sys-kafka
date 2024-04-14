package transport

import (
	"eval-sys-registration-service/cmd/internal/controllers"
	"fmt"

	kafkaWrapper "github.com/engpap/kafka-wrapper-go/pkg"

	"github.com/gin-gonic/gin"
)

func Serve() {
	port := "8081"
	router := initRouter()
	router.Run(":" + port)
	fmt.Println("Server is running on port " + port)
}

func initRouter() *gin.Engine {
	router := gin.Default()

	registrationController := controllers.Controller{}

	go kafkaWrapper.CreateConsumer("course", registrationController.UpdateCourseInMemory)
	go kafkaWrapper.CreateConsumer("project", registrationController.UpdateProjectInMemory)
	go kafkaWrapper.CreateConsumer("grade", registrationController.UpdateGradeInMemory)
	go kafkaWrapper.CreateConsumer("submission", registrationController.UpdateSubmissionInMemory)

	return router
}
