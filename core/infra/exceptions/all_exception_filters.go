package exceptions

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"runtime/debug"
)

type statusCoder interface {
	StatusCode() int
	Error() string
}

type hasDetails interface {
	GetDetails() interface{}
}

func AllExceptionFilter() gin.HandlerFunc {
	isProd := os.Getenv("ENV") == "prod"
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				handleRecoveredError(c, r, isProd)
			}
		}()
		c.Next()
	}
}

func handleRecoveredError(c *gin.Context, r interface{}, isProd bool) {
	status := http.StatusInternalServerError
	message := "Ops! Aconteceu algo de errado do nosso lado. Por favor, entre em contato com o suporte."
	var details interface{} = nil

	err, ok := r.(error)
	if !ok {
		respond(c, status, message, details, isProd)
		return
	}

	if se := extractHttpError(err); se != nil {
		status = se.StatusCode()
		message = se.Error()
	} else {
		message = err.Error()
	}

	if de := extractDetails(err); de != nil {
		details = de
	}

	respond(c, status, message, details, isProd)
}

func extractHttpError(err error) statusCoder {
	var sc statusCoder
	if errors.As(err, &sc) {
		return sc
	}
	return nil
}

func extractDetails(err error) interface{} {
	var d hasDetails
	if errors.As(err, &d) {
		return d.GetDetails()
	}
	return nil
}

func respond(c *gin.Context, status int, message string, details interface{}, isProd bool) {
	resp := gin.H{
		"statusCode": status,
		"message":    message,
		"path":       c.FullPath(),
	}
	if details != nil {
		resp["errors"] = details
	}
	if !isProd {
		resp["stack"] = string(debug.Stack())
	}
	c.AbortWithStatusJSON(status, resp)
}
