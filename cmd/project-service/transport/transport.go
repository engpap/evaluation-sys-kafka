package transport

import (
	kafkaUtils "evaluation-sys-kafka/internal/kafka"
	"evaluation-sys-kafka/pkg/projects/controllers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Serve() {
	port := "8084"
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

	projectController := controllers.Controller{Producer: producer}
	kafkaUtils.SetupCloseProducerHandler(producer)

	go kafkaUtils.CreateConsumer("course", projectController.UpdateCourseInMemory)

	router.POST("/courses/:course-id/projects/create", projectController.CreateProject)
	router.POST("/courses/:course-id/projects/:project-id/submit", projectController.SubmitProjectSolution)
	router.POST("/courses/:course-id/projects/:project-id/submissions/:submission-id/grade", projectController.GradeProjectSolution)

	return router
}
