package transport

import (
	"eval-sys-project-service/cmd/internal/controllers"
	"fmt"

	kafkaWrapper "github.com/engpap/kafka-wrapper-go/pkg"

	"github.com/gin-gonic/gin"
)

func Serve() {
	debug := false
	port := "8080"
	if debug {
		port = "8091"
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

	projectController := controllers.Controller{Producer: producer}
	kafkaWrapper.SetupCloseProducerHandler(producer)

	// consumes on events created by other services
	go kafkaWrapper.CreateConsumer("course", projectController.UpdateCourseInMemory, "project-service")
	go kafkaWrapper.CreateConsumer("enrollment", projectController.UpdateEnrollmentInMemory, "project-service")
	go kafkaWrapper.CreateConsumer("professor", projectController.UpdateProfessorInMemory, "project-service")
	// consumes on events it creates
	go kafkaWrapper.CreateConsumer("project", projectController.UpdateProjectInMemory, "project-service")
	go kafkaWrapper.CreateConsumer("submission", projectController.UpdateSubmissionInMemory, "project-service")
	go kafkaWrapper.CreateConsumer("grade", projectController.UpdateGradeInMemory, "project-service")

	router.POST("/courses/:course-id/projects/create", projectController.CreateProject)                                              // FE done for prof
	router.POST("/courses/:course-id/projects/:project-id/submit", projectController.SubmitProjectSolution)                          // FE done for stud
	router.POST("/courses/:course-id/projects/:project-id/submissions/:submission-id/grade", projectController.GradeProjectSolution) // FE done for prof
	// not required by the specs
	router.GET("/courses/:course-id/projects", projectController.GetCourseProjects)                                                 // FE done for stud
	router.GET("/courses/:course-id/projects/:project-id/submissions", projectController.GetProjectSubmissions)                     // FE done for prof
	router.GET("/courses/:course-id/projects/:project-id/submissions/:submission-id/grades", projectController.GetSubmissionGrades) // FE done for stud

	return router
}
