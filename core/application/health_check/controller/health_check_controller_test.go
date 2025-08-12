package health_check_controller_test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"queue/core/application/health_check/controller"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	r := gin.Default()
	controller := health_check_controller.NewHealthCheckController()
	r.GET("/", controller.HealthCheck)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `Queue on Air!`)
}
