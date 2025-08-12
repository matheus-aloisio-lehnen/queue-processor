package health_check_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthCheckController struct{}

func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

func (h *HealthCheckController) HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "Queue on Air!")
}
