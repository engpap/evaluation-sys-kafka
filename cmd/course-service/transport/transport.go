package transport

import (
	kafkaUtils "evaluation-sys-kafka/internal/kafka"
	"evaluation-sys-kafka/pkg/courses/controllers"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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
	setupCloseHandler(producer)

	router.POST("/courses/create", courseController.CreateCourse)

	return router
}

// Setup clean shutdown on Ctrl+C (SIGINT) or SIGTERM
func setupCloseHandler(producer *kafka.Producer) {
	fmt.Println("Press Ctrl+C to exit.")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		producer.Close()
		fmt.Println("Kafka producer closed.")
		os.Exit(0)
	}()
}
