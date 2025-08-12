package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"queue/core/domain/response"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestData struct {
	Name string `json:"name"`
}

func TestSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/success", func(c *gin.Context) {
		response.Success(c, TestData{Name: "Matheus"}, 200, "Sucesso")
	})

	req, _ := http.NewRequest("GET", "/success", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp response.HttpResponse[TestData]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, "Matheus", resp.Data.Name)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Sucesso", resp.Message)
}

func TestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/error", func(c *gin.Context) {
		response.Error(c, 400, "Erro de validação", []string{"Campo obrigatório"})
	})

	req, _ := http.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var resp response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "Erro de validação", resp.Message)
	assert.ElementsMatch(t, []interface{}{"Campo obrigatório"}, resp.Errors.([]interface{}))
}
