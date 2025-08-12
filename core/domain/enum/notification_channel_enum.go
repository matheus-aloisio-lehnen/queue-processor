package enum

type NotificationChannelEnum string

const (
	EmailChannel            NotificationChannelEnum = "EMAIL"
	SmsChannel              NotificationChannelEnum = "SMS"
	WhatsappChannel         NotificationChannelEnum = "WHATSAPP"
	PushNotificationChannel NotificationChannelEnum = "PUSH_NOTIFICATION"
)
