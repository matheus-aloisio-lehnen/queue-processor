package enum_test

import (
	"github.com/stretchr/testify/assert"
	"queue/core/domain/enum"
	"testing"
)

func TestTopicEnum(t *testing.T) {
	tests := map[string]struct {
		topic    enum.TopicEnum
		expected enum.TopicEnum
	}{
		"Test Notification Topic": {
			topic:    enum.Notification,
			expected: "notifications-sub",
		},
		"Test HmlNotification Topic": {
			topic:    enum.HmlNotification,
			expected: "hml-notifications-sub",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.topic)
		})
	}
}
