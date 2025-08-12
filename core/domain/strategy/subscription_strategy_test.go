package strategy_test

import (
	"github.com/stretchr/testify/assert"
	"queue/core/domain/enum"
	"queue/core/domain/strategy"
	"queue/core/domain/types"
	"testing"
)

func TestNewSubscriptionHandlerStrategy_Prod(t *testing.T) {
	cfg := &types.Config{Environment: "prod"}
	handlerFactory := strategy.NewSubscriptionHandlerStrategy(cfg)

	handler := handlerFactory.GetHandler(string(enum.Notification))
	assert.NotNil(t, handler, "deveria retornar um handler para topic de produção")
	assert.Nil(t, handlerFactory.GetHandler(string(enum.HmlNotification)), "não deveria retornar handler para hml em produção")
}

func TestNewSubscriptionHandlerStrategy_Hml(t *testing.T) {
	cfg := &types.Config{Environment: "hml"}
	handlerFactory := strategy.NewSubscriptionHandlerStrategy(cfg)

	handler := handlerFactory.GetHandler(string(enum.HmlNotification))
	assert.NotNil(t, handler, "deveria retornar um handler para topic de hml")
	assert.Nil(t, handlerFactory.GetHandler(string(enum.Notification)), "não deveria retornar handler para prod em hml")
}

func TestNewSubscriptionHandlerStrategy_Empty(t *testing.T) {
	cfg := &types.Config{Environment: "other"}
	handlerFactory := strategy.NewSubscriptionHandlerStrategy(cfg)

	assert.Nil(t, handlerFactory.GetHandler(string(enum.Notification)), "não deveria retornar handler para nenhum topic")
	assert.Nil(t, handlerFactory.GetHandler(string(enum.HmlNotification)), "não deveria retornar handler para nenhum topic")
}
