package config

type URLsConfig struct {
	UserServiceURL    string
	ProjectServiceURL string
	CourseServiceURL  string
}

const (
	CourseServiceURL  = "http://Eval-sys-course-env-docker.eba-ij33i5hc.eu-north-1.elasticbeanstalk.com"
	ProjectServiceURL = "http://Eval-sys-project-env-docker.eba-tx9pz6g2.eu-north-1.elasticbeanstalk.com"
	UserServiceURL    = "http://Eval-sys-user-env-docker.eba-qj3fh5wc.eu-north-1.elasticbeanstalk.com"
)

const (
	DebugCourseServiceURL  = "http://localhost:8090"
	DebugProjectServiceURL = "http://localhost:8091"
	// 8092 is for registration-service, which is not recheaded by the cli
	DebugUserServiceURL = "http://localhost:8093"
)

const debug = true

func init() {
	if debug {
		URLs = URLsConfig{
			UserServiceURL:    DebugUserServiceURL,
			ProjectServiceURL: DebugProjectServiceURL,
			CourseServiceURL:  DebugCourseServiceURL,
		}
	} else {
		URLs = URLsConfig{
			UserServiceURL:    UserServiceURL,
			ProjectServiceURL: ProjectServiceURL,
			CourseServiceURL:  CourseServiceURL,
		}
	}

}

var URLs URLsConfig
