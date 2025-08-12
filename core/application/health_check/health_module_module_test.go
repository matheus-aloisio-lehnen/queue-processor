package healthcheckmodule_test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"queue/core/application/health_check"
	"testing"
)

func TestHealthCheckModule_RegisterRoutes(t *testing.T) {
	r := gin.Default()
	module := healthcheckmodule.NewHealthCheckModule()
	module.RegisterRoutes(r)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody := `Queue on Air!`
	assert.Equal(t, expectedBody, w.Body.String(), "A resposta não corresponde à esperada")
}
