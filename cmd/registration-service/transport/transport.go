package transport

import (
	"evaluation-sys-kafka/cmd/registration-service/controllers"
	kafkaUtils "evaluation-sys-kafka/internal/kafkawrapper"
	"fmt"

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

	go kafkaUtils.CreateConsumer("course", registrationController.UpdateCourseInMemory)
	go kafkaUtils.CreateConsumer("project", registrationController.UpdateProjectInMemory)
	go kafkaUtils.CreateConsumer("grade", registrationController.UpdateGradeInMemory)
	go kafkaUtils.CreateConsumer("submission", registrationController.UpdateSubmissionInMemory)

	return router
}
