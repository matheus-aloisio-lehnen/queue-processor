package classtransformer_test

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"queue/core/infra/middleware"
	"testing"
)

func TestUseClassTransformerMiddleware(t *testing.T) {
	type TestDTO struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"required"`
	}
	r := gin.Default()
	r.POST("/test", classtransformer.UseClassTransformerMiddleware(TestDTO{}), func(c *gin.Context) {
		dto, exists := c.Get("dto")
		assert.True(t, exists)
		assert.IsType(t, &TestDTO{}, dto)
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer([]byte(`{"name":"John", "age":30}`)))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Status esperado: 200")
}

func TestClassTransformerMiddleware_WithInvalidDTO(t *testing.T) {
	type TestDTO struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"required"`
	}
	r := gin.Default()
	r.POST("/test", classtransformer.UseClassTransformerMiddleware(&TestDTO{}), func(c *gin.Context) {
		// Não será chamado, pois o middleware barrará aqui
		t.FailNow()
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer([]byte(`{"name":"", "age":0}`)))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code, "Status esperado: 400")
}
