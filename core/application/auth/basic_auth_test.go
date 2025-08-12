package auth_test

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"queue/core/application/auth"
	"testing"
)

func TestBasicAuthMiddleware(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	validUsername := "user"
	validPassword := "pass"

	// Middleware sendo testado
	router.GET("/protected", auth.BasicAuthMiddleware(validUsername, validPassword), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test 1: Sucesso - Credenciais corretas
	t.Run("Success - Valid credentials", func(t *testing.T) {
		authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(validUsername+":"+validPassword))
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", authHeader)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var responseBody map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Equal(t, "success", responseBody["message"])
	})

	// Test 2: Erro - Header ausente
	t.Run("Error - Missing Authorization header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var responseBody map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Equal(t, "Autenticação de headers requerida", responseBody["message"])
	})

	// Test 3: Erro - Formato de header inválido
	t.Run("Error - Invalid Authorization header format", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var responseBody map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Equal(t, "Autenticação de headers inválida", responseBody["message"])
	})

	// Test 4: Erro - Falha ao decodificar header
	t.Run("Error - Failed to decode header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Basic invalidbase64==")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var responseBody map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Equal(t, "Falha ao decodificar o header", responseBody["message"])
	})

	// Test 5: Erro - Credenciais inválidas
	t.Run("Error - Invalid credentials", func(t *testing.T) {
		invalidCredentials := "Basic " + base64.StdEncoding.EncodeToString([]byte("invalidUser:invalidPass"))
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", invalidCredentials)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var responseBody map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Equal(t, "Credenciais inválidas", responseBody["message"])
	})
}
