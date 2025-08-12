package enum_test

import (
	"github.com/stretchr/testify/assert"
	"queue/core/domain/enum"
	"testing"
)

func TestNotificationChannelEnum(t *testing.T) {
	tests := map[string]struct {
		channel     enum.NotificationChannelEnum
		expected    string
		expectedErr bool
	}{
		"Email Channel": {
			channel:  enum.EmailChannel,
			expected: "EMAIL",
		},
		"Ss Channel": {
			channel:  enum.SmsChannel,
			expected: "SMS",
		},
		"Whatsapp Channel": {
			channel:  enum.WhatsappChannel,
			expected: "WHATSAPP",
		},
		"Push Notification Channel": {
			channel:  enum.PushNotificationChannel,
			expected: "PUSH_NOTIFICATION",
		},
		"Invalid Email Channel Comparison": {
			channel:     enum.EmailChannel,
			expected:    "SMS", // Valor incorreto, deve falhar
			expectedErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.expectedErr {
				assert.NotEqual(t, string(tt.channel), tt.expected, "O canal não deve ser igual ao esperado")
			} else {
				assert.Equal(t, string(tt.channel), tt.expected, "O canal deve ser igual ao valor esperado")
			}
		})
	}

	t.Run("Testar tipo NotificationChannelEnum", func(t *testing.T) {
		var channel enum.NotificationChannelEnum
		assert.IsType(t, string(""), string(channel), "NotificationChannelEnum não é do tipo string")
	})
}
