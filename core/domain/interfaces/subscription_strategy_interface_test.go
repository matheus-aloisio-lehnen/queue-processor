package interfaces_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"queue/core/domain/interfaces"
)

// MockHandlerStrategy implementa HandlerStrategyInterface para teste
type MockHandlerStrategy struct {
	Handlers map[string]interfaces.ISubscribeHandler
}

func (m *MockHandlerStrategy) GetHandler(topic string) interfaces.ISubscribeHandler {
	return m.Handlers[topic]
}

type DummyHandler struct{}

func (d *DummyHandler) Handle(sub interfaces.ISubscription) error {
	return nil
}

func TestHandlerStrategyInterface_GetHandler(t *testing.T) {
	handler1 := &DummyHandler{}
	handler2 := &DummyHandler{}

	strategy := &MockHandlerStrategy{
		Handlers: map[string]interfaces.ISubscribeHandler{
			"topic-1": handler1,
			"topic-2": handler2,
		},
	}

	t.Run("Retorna handler correto para tópico conhecido", func(t *testing.T) {
		assert.Equal(t, handler1, strategy.GetHandler("topic-1"))
		assert.Equal(t, handler2, strategy.GetHandler("topic-2"))
	})

	t.Run("Retorna nil para tópico desconhecido", func(t *testing.T) {
		assert.Nil(t, strategy.GetHandler("topic-xyz"))
	})
}
