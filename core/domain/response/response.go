package response

import "github.com/gin-gonic/gin"

type HttpResponse[T any] struct {
	Data       T      `json:"data"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type IList[T any] struct {
	List      []T   `json:"list"`
	TotalRows int64 `json:"totalRows"`
}

type ErrorResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Errors     interface{} `json:"errors,omitempty"`
}

func Success[T any](c *gin.Context, data T, statusCode int, message string) {
	c.JSON(statusCode, HttpResponse[T]{
		Data:       data,
		StatusCode: statusCode,
		Message:    message,
	})
}

func Error(c *gin.Context, statusCode int, message string, details interface{}) {
	c.AbortWithStatusJSON(statusCode, ErrorResponse{
		StatusCode: statusCode,
		Message:    message,
		Errors:     details,
	})

}
