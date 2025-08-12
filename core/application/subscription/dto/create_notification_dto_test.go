package subscriptiondto_test

import (
	"queue/core/application/subscription/dto"
	"queue/core/domain/enum"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var validate = validator.New()

func TestCreateNotificationDto_Validacoes(t *testing.T) {
	tests := map[string]struct {
		dto        subscriptiondto.CreateNotificationDto
		esperaErro bool
		camposErro []string
	}{
		"Payload de email válido": {
			dto: subscriptiondto.CreateNotificationDto{
				UserID:    1,
				UserName:  "João",
				Channel:   enum.EmailChannel,
				Subject:   "Boas-vindas",
				Recipient: "joao@email.com",
				Payload: subscriptiondto.EmailPayload{
					HTML:      "<h1>Olá</h1>",
					PlainText: "Olá",
				},
			},
		},
		"Payload de mensagem (SMS) válido": {
			dto: subscriptiondto.CreateNotificationDto{
				UserID:    2,
				UserName:  "Maria",
				Channel:   enum.SmsChannel,
				Recipient: "+5511999999999",
				Payload: subscriptiondto.MessagePayload{
					Message: "Olá, tudo bem?",
				},
			},
		},
		"Payload de push notification válido": {
			dto: subscriptiondto.CreateNotificationDto{
				UserID:    3,
				UserName:  "Pedro",
				Channel:   enum.PushNotificationChannel,
				Recipient: "device-token",
				Payload: subscriptiondto.PushNotificationPayload{
					DeviceToken: "device-token",
					Title:       "Alerta",
					Body:        "Nova mensagem disponível",
					Data: map[string]any{
						"foo": "bar",
					},
				},
			},
		},
		"DTO inválido: campos obrigatórios ausentes": {
			dto:        subscriptiondto.CreateNotificationDto{},
			esperaErro: true,
			camposErro: []string{"UserName", "Recipient", "Payload"},
		},
	}

	for nome, tc := range tests {
		t.Run(nome, func(t *testing.T) {
			err := validate.Struct(tc.dto)

			if tc.esperaErro {
				assert.Error(t, err)
				for _, campo := range tc.camposErro {
					assert.Contains(t, err.Error(), campo)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
