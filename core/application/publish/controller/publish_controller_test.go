package publishcontroller_test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"queue/core/application/publish/controller"
	"queue/core/application/publish/dto"
	"queue/core/application/publish/service"
	"queue/core/infra/mock"
	"testing"
)

// Configuração simplificada do router, "injeção" manual no contexto da request do "dto".
func setupRouter(srv publishservice.PublishService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	ctrl := publishcontroller.NewPublishController(srv)
	r.POST("/publish", func(c *gin.Context) {
		var input publishdto.InputDto
		body, _ := io.ReadAll(c.Request.Body)
		_ = json.Unmarshal(body, &input)
		c.Set("dto", &input)
		ctrl.Publish(c)
	})
	return r
}

func TestPublishController(t *testing.T) {
	type testCase struct {
		name        string
		topicName   string
		payload     any
		wantStatus  int
		wantContent string
	}

	testCases := []testCase{
		{
			name:      "sucesso",
			topicName: "ok", // não é "fail", então retorna sucesso
			payload: publishdto.InputDto{
				Meta: publishdto.MetaDto{Topic: "ok"},
				Data: map[string]string{"foo": "bar"},
			},
			wantStatus:  http.StatusOK,
			wantContent: "Mensagem publicada com sucesso",
		},
		{
			name:      "erro no publish",
			topicName: "fail",
			payload: publishdto.InputDto{
				Meta: publishdto.MetaDto{Topic: "fail"},
				Data: map[string]string{"x": "y"},
			},
			wantStatus:  http.StatusBadRequest,
			wantContent: "erro simulado",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Usa MockPubSubClientAdapter ajustado
			mockClient := mock.NewMockPubSubClientAdapter()
			svc, _ := publishservice.NewPublishService(mockClient)
			router := setupRouter(*svc)

			body, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/publish", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.wantStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.wantContent)
		})
	}

	// Caso DTO não presente
	t.Run("dto não presente", func(t *testing.T) {
		mockClient := mock.NewMockPubSubClientAdapter()
		svc, _ := publishservice.NewPublishService(mockClient)
		ctrl := publishcontroller.NewPublishController(*svc)
		router := gin.Default()
		router.POST("/publish", func(c *gin.Context) {
			ctrl.Publish(c)
		})
		req, _ := http.NewRequest(http.MethodPost, "/publish", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Formato de mensagem inválido")
	})
}
