package handlers

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"queue/core/application/subscription/dto"
	"queue/core/domain/enum"
	"queue/core/domain/interfaces"
	"queue/core/domain/types"
	"queue/core/infra/http_service"
)

type NotificationHandler struct {
	Config      *types.Config
	httpService *httpservice.HttpService
}

func NewNotificationHandler(cfg *types.Config) *NotificationHandler {
	return &NotificationHandler{
		Config:      cfg,
		httpService: httpservice.NewHttpService(),
	}
}

func (h *NotificationHandler) Handle(sub interfaces.ISubscription) error {
	ctx := context.Background()
	return sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		dto, err := parseNotificationMessage(msg.Data)
		if err != nil {
			msg.Nack()
			return
		}
		if err := h.httpService.SendNotification(dto, h.Config); err != nil {
			msg.Nack()
			return
		}
		msg.Ack()
	})
}

func parseNotificationMessage(body []byte) (*subscriptiondto.CreateNotificationDto, error) {
	// Extrai campo "data"
	data, err := extractDataField(body)
	if err != nil {
		return nil, err
	}

	// Converte para DTO
	dto, err := parseNotificationData(data)
	if err != nil {
		return nil, err
	}

	return dto, nil
}

func extractDataField(raw []byte) ([]byte, error) {
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("erro ao decodificar envelope: %w", err)
	}
	data, ok := envelope["data"]
	if !ok {
		return nil, errors.New("campo 'data' não encontrado na mensagem")
	}
	return data, nil
}

func parseNotificationData(data []byte) (*subscriptiondto.CreateNotificationDto, error) {
	var dto subscriptiondto.CreateNotificationDto
	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, fmt.Errorf("erro ao decodificar dados: %w", err)
	}

	if dto.UserID == 0 || dto.Channel == "" || dto.Recipient == "" {
		return nil, errors.New("dados obrigatórios ausentes")
	}

	dto.Payload = parseTypedPayload(dto.Channel, dto.Payload)
	return &dto, nil
}

func parseTypedPayload(channel enum.NotificationChannelEnum, raw any) any {
	switch channel {
	case enum.EmailChannel:
		var p subscriptiondto.EmailPayload
		_ = remarshal(raw, &p)
		return p
	case enum.SmsChannel, enum.WhatsappChannel:
		var p subscriptiondto.MessagePayload
		_ = remarshal(raw, &p)
		return p
	case enum.PushNotificationChannel:
		var p subscriptiondto.PushNotificationPayload
		_ = remarshal(raw, &p)
		return p
	default:
		return raw
	}
}

func remarshal(from any, to any) error {
	b, err := json.Marshal(from) // transforma map[any] em []byte
	if err != nil {
		return err
	}
	return json.Unmarshal(b, to) // agora transforma []byte na struct correta
}
