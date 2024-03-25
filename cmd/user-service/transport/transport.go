package transport

import (
	kafkaUtils "evaluation-sys-kafka/internal/kafka"
	"evaluation-sys-kafka/pkg/users/controllers"
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

	//producer, err := kafkaUtils.CreateProducer()
	//if err != nil {
	//	panic(err)
	//}
	//userController := controllers.Controller{Producer: producer}
	//afkaUtils.SetupCloseHandler(producer)

	// TODO: set up consumer

	//kafkaUtils.SetupCloseConsumerHandler(consumer)

	userController := controllers.Controller{}

	go kafkaUtils.CreateConsumer("course", &userController.ConsumerOutput)

	router.GET("/student/get-course", userController.GetCourses)

	return router
}
