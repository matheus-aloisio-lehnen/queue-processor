package publishdto

import (
	"github.com/go-playground/validator/v10"
	"queue/core/infra/validation"
)

var createInputErrors = map[string]string{
	"Meta.Topic.required": "O tópico é obrigatório.",
	"Data.required":       "Os dados são obrigatórios.",
}

type MetaDto struct {
	AccessToken *string `json:"access_token,omitempty" validate:"omitempty"`
	Topic       string  `json:"topic" validate:"required"`
}

type InputDto struct {
	Meta MetaDto     `json:"meta" validate:"required"`
	Data interface{} `json:"data" validate:"required"`
}

func (dto *InputDto) ValidationMessages(ve validator.ValidationErrors) map[string]string {
	return validation.ValidationMessages(ve, createInputErrors)
}
