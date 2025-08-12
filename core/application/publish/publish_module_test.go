package publishmodule_test

import (
	"context"
	"queue/core/application/publish"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"queue/core/domain/interfaces"
)

type MockPubSubClient struct{}

func (m *MockPubSubClient) Subscriptions(ctx context.Context) interfaces.ISubscriptionIterator {
	return nil
}
func (m *MockPubSubClient) Topic(id string) interfaces.ITopic { return nil }
func (m *MockPubSubClient) Close() error                      { return nil }

func TestNewPublishModule_Success(t *testing.T) {
	client := &MockPubSubClient{}
	module, err := publishmodule.NewPublishModule(client)
	assert.NoError(t, err)
	assert.NotNil(t, module)
	assert.NotNil(t, module.Controller)
}

func TestPublishModule_RegisterRoutes(t *testing.T) {
	client := &MockPubSubClient{}
	module, err := publishmodule.NewPublishModule(client)
	assert.NoError(t, err)
	router := gin.Default()
	module.RegisterRoutes(router)
	req := httptest.NewRequest(http.MethodPost, "/publish", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.NotEqual(t, http.StatusNotFound, w.Code)
	assert.Contains(t, []int{http.StatusBadRequest, http.StatusInternalServerError, http.StatusOK}, w.Code)
}
