package auth

import (
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"queue/core/domain/response"
	"strings"
)

func BasicAuthMiddleware(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/" {
			c.Next()
			return
		}
		if err := authenticate(c, username, password); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		c.Next()
	}
}

func authenticate(c *gin.Context, validUsername, validPassword string) error {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return errors.New("Autenticação de headers requerida")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		return errors.New("Autenticação de headers inválida")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return errors.New("Falha ao decodificar o header")
	}

	credentials := strings.Split(string(decoded), ":")
	if len(credentials) != 2 || credentials[0] != validUsername || credentials[1] != validPassword {
		return errors.New("Credenciais inválidas")
	}

	return nil
}
