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
	courseController := controllers.Controller{Producer: producer}
	kafkaUtils.SetupCloseProducerHandler(producer)

	router.POST("/courses/create", courseController.CreateCourse)

	return router
}
