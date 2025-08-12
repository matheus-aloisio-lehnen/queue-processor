package strategy

import (
	"queue/core/domain/enum"
	"queue/core/domain/interfaces"
	"queue/core/domain/strategy/handlers"
	"queue/core/domain/types"
)

type SubscriptionHandlerStrategy struct {
	handler map[enum.TopicEnum]interfaces.ISubscribeHandler
}

func NewSubscriptionHandlerStrategy(cfg *types.Config) *SubscriptionHandlerStrategy {
	handler := map[enum.TopicEnum]interfaces.ISubscribeHandler{}
	if cfg.Environment == "prod" {
		handler[enum.Notification] = handlers.NewNotificationHandler(cfg)
	}
	if cfg.Environment == "hml" {
		handler[enum.HmlNotification] = handlers.NewNotificationHandler(cfg)
	}
	return &SubscriptionHandlerStrategy{handler: handler}
}

func (f *SubscriptionHandlerStrategy) GetHandler(topic string) interfaces.ISubscribeHandler {
	if handler, ok := f.handler[enum.TopicEnum(topic)]; ok {
		return handler
	}
	return nil
}
