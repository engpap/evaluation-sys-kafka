Current flow:
- Admin post a new course through Course Service. Course Service is the producer.
- Student Service is the conusmer for the topic course. Course's data is replicated from the Course Service to the Student Service, where it can be queried locally. (see pag. 37)


// Student Routes
router.POST("/courses/:course-id/enroll", controllers.EnrollCourse)
router.POST("/projects/:project-id/submit", controllers.SubmitProject)
router.GET("/projects/my-submissions", controllers.GetMySubmissions)

// Admin Routes
router.POST("/courses/:course-id/delete", controllers.DeleteCourse)
router.POST("/courses/create", controllers.CreateCourse)

// Professor Routes
router.POST("/projects/:project-id/grade", controllers.GradeProject)
router.POST("/projects/create", controllers.CreateProject)

| Microservice | Generated Event | Listened Event | Description |
| --- | --- | --- | --- |
| User Service | CourseEnrollment | CourseCompleted | manages personal data of registered students and professors; |
| Course Service | 
 | CourseEnrollment
ProjectCreated | manages courses, which consists of several projects; |
| Project Service | ProjectCreated
SolutionSubmitted
SolutionGraded |  | manages submission and grading of individual projects; |
| Registration Service | CourseCompleted | SolutionGraded | handles registration of grades for completed courses. |

| Actors | Actions |
| --- | --- |
| Students | enroll for courses, submit their solutions for projects, and check the status of their submissions; |
| Admins |  add new courses and remove old courses; |
| Professors | post projects and grade the solutions submitted by students |