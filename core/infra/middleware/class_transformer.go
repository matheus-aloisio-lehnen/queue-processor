package classtransformer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	validationerrormessageprovider "queue/core/domain/interfaces/validation_error_message_provider"
	"queue/core/domain/response"
	"queue/core/domain/structs"
	"queue/core/infra/utils/pipe"
	"reflect"
)

var validate = validator.New()

func UseClassTransformerMiddleware(dtoType interface{}) gin.HandlerFunc {
	typ := reflect.TypeOf(dtoType)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return func(c *gin.Context) {
		c.Set("dtoType", typ)
		ClassTransformerMiddleware()(c)
	}
}

func ClassTransformerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawType, exists := c.Get("dtoType")
		if !exists {
			handleError(c, "Tipo de DTO não configurado", fmt.Errorf("Tipo de DTO não encontrado no contexto"))
			return
		}
		typ, ok := rawType.(reflect.Type)
		if !ok {
			handleError(c, "Tipo de DTO inválido", fmt.Errorf("Tipo assert error"))
			return
		}
		body, err := captureRequestBody(c)
		if err != nil {
			handleError(c, "Erro ao ler corpo da requisição", err)
			return
		}
		dto := reflect.New(typ).Interface()
		err = BindAndValidate(body, dto)
		if err != nil {
			handleError(c, err.Error(), err)
			return
		}
		c.Set("dto", dto)
		c.Next()
	}
}

func captureRequestBody(c *gin.Context) ([]byte, error) {
	raw, err := c.GetRawData()
	if err != nil {
		return nil, fmt.Errorf("Erro ao ler corpo da requisição: %v", err)
	}
	return raw, nil
}

func BindAndValidate(data []byte, dto interface{}) error {
	err := decodeJSON(data, dto)
	if err != nil {
		return err
	}
	err = validate.Struct(dto)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			if provider, ok := dto.(validationerrormessageprovider.ValidationMessageProvider); ok {
				msgsMap := provider.ValidationMessages(ve)
				msgs := make([]string, 0, len(msgsMap))
				for _, v := range msgsMap {
					msgs = append(msgs, v)
				}
				return &structs.ValidationMessagesError{Messages: msgs}

			}
			return &structs.ValidationMessagesError{Messages: []string{"Erro de validação"}}
		}
	}
	pipe.TrimPipe(dto)
	return nil
}

func decodeJSON(data []byte, out interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("Erro ao decodificar dados: %v", err)
	}
	return nil
}

func handleError(c *gin.Context, errorMessage string, err error) {
	details := fmt.Sprintf("%s: %v", errorMessage, err)
	response.Error(c, http.StatusBadRequest, errorMessage, details)
	return
}
