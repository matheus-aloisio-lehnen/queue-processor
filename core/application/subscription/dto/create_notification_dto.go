package subscriptiondto

import "queue/core/domain/enum"

type EmailPayload struct {
	HTML      string `json:"html"`
	PlainText string `json:"plainText"`
}

type MessagePayload struct {
	Message string `json:"message"`
}

type PushNotificationPayload struct {
	DeviceToken string         `json:"deviceToken"`
	Title       string         `json:"title"`
	Body        string         `json:"body"`
	Data        map[string]any `json:"data,omitempty"`
}

type CreateNotificationDto struct {
	UserID    int64                        `json:"userId" validate:"required"`
	UserName  string                       `json:"userName" validate:"required"`
	Channel   enum.NotificationChannelEnum `json:"channel" validate:"required"`
	Subject   string                       `json:"subject,omitempty"`
	Recipient string                       `json:"recipient" validate:"required"`
	Payload   any                          `json:"payload" validate:"required"`
}
