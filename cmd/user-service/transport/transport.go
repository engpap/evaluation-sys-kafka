package transport

import (
	kafkaUtils "evaluation-sys-kafka/internal/kafka"
	"evaluation-sys-kafka/pkg/users/controllers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Serve() {
	port := "8083"
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
	userController := controllers.Controller{Producer: producer}
	kafkaUtils.SetupCloseProducerHandler(producer)
	// TODO: set up consumer handler

	router.POST("/users/stud/create", userController.CreateStudent)
	router.POST("/users/prof/create", userController.CreateProfessor)

	return router
}
