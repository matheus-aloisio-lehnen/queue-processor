package exceptions_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"queue/core/infra/exceptions"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Tipos de erro personalizados para testar os diferentes cenários
type simpleError struct{}

func (e *simpleError) Error() string { return "erro simples" }

type codedError struct{}

func (e *codedError) Error() string   { return "erro com status" }
func (e *codedError) StatusCode() int { return 418 }

type detailedError struct{}

func (e *detailedError) Error() string   { return "erro com detalhes" }
func (e *detailedError) StatusCode() int { return 422 }
func (e *detailedError) GetDetails() any {
	return []string{"campo X obrigatório", "campo Y inválido"}
}

func TestAllExceptionFilter_WithCodedError(t *testing.T) {
	// Ambiente de teste
	t.Setenv("ENV", "dev")
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(exceptions.AllExceptionFilter()) // Usando o filtro de exceções
	router.GET("/test", func(c *gin.Context) {
		panic(&codedError{}) // Simulando um panic com um erro personalizado
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	// Verificando se o código de status e a mensagem estão corretos
	assert.Equal(t, 418, w.Code)
	assert.Contains(t, w.Body.String(), "erro com status")
	assert.Contains(t, w.Body.String(), "statusCode")
	assert.Contains(t, w.Body.String(), "path")
}

func TestAllExceptionFilter_WithDetailedError(t *testing.T) {
	// Ambiente de teste
	t.Setenv("ENV", "dev")
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(exceptions.AllExceptionFilter()) // Usando o filtro de exceções
	router.GET("/details", func(c *gin.Context) {
		panic(&detailedError{}) // Simulando um panic com erro detalhado
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/details", nil)
	router.ServeHTTP(w, req)

	// Verificando se o código de status, mensagem e detalhes estão corretos
	assert.Equal(t, 422, w.Code)
	assert.Contains(t, w.Body.String(), "erro com detalhes")
	assert.Contains(t, w.Body.String(), "campo X obrigatório")
	assert.Contains(t, w.Body.String(), "campo Y inválido")
}

func TestAllExceptionFilter_WithBasicError(t *testing.T) {
	// Ambiente de teste
	t.Setenv("ENV", "dev")
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(exceptions.AllExceptionFilter()) // Usando o filtro de exceções
	router.GET("/basic", func(c *gin.Context) {
		panic(errors.New("erro cru")) // Simulando um erro genérico
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/basic", nil)
	router.ServeHTTP(w, req)

	// Verificando se o código de status é 500 e a mensagem do erro genérico
	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "erro cru")
}

func TestAllExceptionFilter_WithNonErrorPanic(t *testing.T) {
	// Ambiente de teste
	t.Setenv("ENV", "dev")
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(exceptions.AllExceptionFilter()) // Usando o filtro de exceções
	router.GET("/string", func(c *gin.Context) {
		panic("texto bruto") // Simulando um panic com texto simples
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/string", nil)
	router.ServeHTTP(w, req)

	// Verificando se o código de status é 500 e a mensagem padrão
	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "Ops! Aconteceu algo de errado")
}
