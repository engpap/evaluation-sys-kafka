package transport

import (
	"eval-sys-user-service/cmd/internal/controllers"
	"fmt"

	kafkaWrapper "github.com/engpap/kafka-wrapper-go/pkg"

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

	producer, err := kafkaWrapper.CreateProducer()
	if err != nil {
		panic(err)
	}
	userController := controllers.Controller{Producer: producer}
	kafkaWrapper.SetupCloseProducerHandler(producer)
	// TODO: set up consumer handler

	router.POST("/users/student/create", userController.CreateStudent)     // FE done for admin
	router.POST("/users/professor/create", userController.CreateProfessor) // FE done for admin

	return router
}
