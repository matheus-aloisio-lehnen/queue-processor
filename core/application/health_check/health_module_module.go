package healthcheckmodule

import (
	"github.com/gin-gonic/gin"
	"queue/core/application/health_check/controller"
)

type HealthCheckModule struct {
	Controller *health_check_controller.HealthCheckController
}

func NewHealthCheckModule() *HealthCheckModule {
	ctrl := health_check_controller.NewHealthCheckController()
	return &HealthCheckModule{
		Controller: ctrl,
	}
}

func (m *HealthCheckModule) RegisterRoutes(router *gin.Engine) {
	router.GET("/", m.Controller.HealthCheck)
}
