package publishdto_test

import (
	"queue/core/application/publish/dto"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var validate = validator.New()

func TestInputDto_DeveSerValidoQuandoTodosOsCamposPresentes(t *testing.T) {
	input := publishdto.InputDto{
		Meta: publishdto.MetaDto{
			Topic: "meu-topico",
		},
		Data: map[string]interface{}{
			"mensagem": "teste",
		},
	}

	err := validate.Struct(input)
	assert.NoError(t, err)
}

func TestInputDto_DeveRetornarErroQuandoTopicForVazio(t *testing.T) {
	input := publishdto.InputDto{
		Meta: publishdto.MetaDto{
			Topic: "",
		},
		Data: map[string]interface{}{
			"mensagem": "teste",
		},
	}

	err := validate.Struct(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Topic")
}

func TestInputDto_DeveRetornarErroQuandoMetaForVazio(t *testing.T) {
	input := publishdto.InputDto{
		Data: map[string]interface{}{
			"mensagem": "teste",
		},
	}

	err := validate.Struct(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Meta")
}

func TestInputDto_DeveRetornarErroQuandoDataForNil(t *testing.T) {
	input := publishdto.InputDto{
		Meta: publishdto.MetaDto{
			Topic: "meu-topico",
		},
		Data: nil,
	}

	err := validate.Struct(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Data")
}
