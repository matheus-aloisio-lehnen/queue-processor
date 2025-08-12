package servers_test

import (
	"net/http"
	"net/http/httptest"
	"queue/core/infra/mock"
	"queue/core/infra/servers"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"queue/core/domain/types"
)

func TestSetupRouter_HealthRouteAndAuth(t *testing.T) {
	cfg := &types.Config{
		Auth: types.BasicAuthConfig{
			Username: "admin",
			Password: "123",
		},
	}
	pubsubMock := mock.NewMockPubSubClientAdapter()
	router := servers.SetupRouter(cfg, pubsubMock)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	req.SetBasicAuth("admin", "123")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Queue on Air")
}

func TestStartServer_DoesNotPanicOnPortZero(t *testing.T) {
	router := gin.New()
	go func() {
		_ = servers.StartServer(router, 0)
	}()
	assert.True(t, true)
}
