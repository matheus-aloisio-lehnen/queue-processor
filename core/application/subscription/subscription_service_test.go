package subscriptionservice_test

import (
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"queue/core/domain/interfaces"
	"queue/core/domain/types"
	"queue/core/infra/mock"
	"testing"

	subscriptionservice "queue/core/application/subscription"
)

// Handler falso para simular comportamentos nos testes
type fakeHandler struct {
	Called   bool
	CallWith interfaces.ISubscription
	Err      error
}

func (f *fakeHandler) Handle(sub interfaces.ISubscription) error {
	f.Called = true
	f.CallWith = sub
	return f.Err
}

// Mock robusto de HandlerStrategyInterface
type fakeHandlerStrategy struct {
	handlers map[string]interfaces.ISubscribeHandler
}

func (f *fakeHandlerStrategy) GetHandler(topic string) interfaces.ISubscribeHandler {
	return f.handlers[topic]
}

func TestSubscriptionService_Listen_WithHandler(t *testing.T) {
	cfg := &types.Config{}
	handler := &fakeHandler{}
	handlerStrategy := &fakeHandlerStrategy{
		handlers: map[string]interfaces.ISubscribeHandler{
			"sub-ok": handler,
		},
	}
	client := mock.NewMockPubSubClientAdapter("sub-ok")
	service := subscriptionservice.NewSubscriptionService(client, cfg)
	service.HandlerStrategy = handlerStrategy // sobrescrevendo após construção

	err := service.Listen(context.Background())
	if err != nil {
		t.Fatalf("Listen retornou erro inesperado: %v", err)
	}
	if !handler.Called {
		t.Errorf("Esperava que handler fosse chamado para subscription existente.")
	}
}

func TestSubscriptionService_Listen_WithoutHandler(t *testing.T) {
	cfg := &types.Config{}
	handlerStrategy := &fakeHandlerStrategy{
		handlers: map[string]interfaces.ISubscribeHandler{},
	}
	client := mock.NewMockPubSubClientAdapter("sub-sem-handler")
	service := subscriptionservice.NewSubscriptionService(client, cfg)
	service.HandlerStrategy = handlerStrategy

	err := service.Listen(context.Background())
	if err != nil {
		t.Fatalf("Listen sem handler deveria terminar sem erro, retornou: %v", err)
	}
}

func TestSubscriptionService_Listen_HandlerError(t *testing.T) {
	cfg := &types.Config{}
	handler := &fakeHandler{Err: errors.New("erro simulado")}
	handlerStrategy := &fakeHandlerStrategy{
		handlers: map[string]interfaces.ISubscribeHandler{
			"sub-err": handler,
		},
	}
	client := mock.NewMockPubSubClientAdapter("sub-err")
	service := subscriptionservice.NewSubscriptionService(client, cfg)
	service.HandlerStrategy = handlerStrategy

	err := service.Listen(context.Background())
	if err == nil {
		t.Fatal("Listen deveria retornar erro quando handler falha.")
	}
}

type iterErr struct{}

func (e *iterErr) Next() (interfaces.ISubscription, error) {
	return nil, errors.New("falha iterator")
}

type clientErr struct{}

func (c *clientErr) Subscriptions(_ context.Context) interfaces.ISubscriptionIterator {
	return &iterErr{}
}
func (c *clientErr) Topic(id string) interfaces.ITopic { return nil }
func (c *clientErr) Close() error                      { return nil }

func TestSubscriptionService_Listen_IteratorError(t *testing.T) {
	cfg := &types.Config{}
	handlerStrategy := &fakeHandlerStrategy{handlers: map[string]interfaces.ISubscribeHandler{}}
	service := subscriptionservice.NewSubscriptionService(&clientErr{}, cfg)
	service.HandlerStrategy = handlerStrategy

	err := service.Listen(context.Background())
	if err == nil || err.Error() != "falha iterator" {
		t.Fatalf("Listen deveria repassar erro inesperado do iterator")
	}
}

type doneIter struct{}

func (d *doneIter) Next() (interfaces.ISubscription, error) {
	return nil, iterator.Done
}

type doneClient struct{}

func (c *doneClient) Subscriptions(_ context.Context) interfaces.ISubscriptionIterator {
	return &doneIter{}
}
func (c *doneClient) Topic(id string) interfaces.ITopic { return nil }
func (c *doneClient) Close() error                      { return nil }

func TestSubscriptionService_Listen_Empty(t *testing.T) {
	cfg := &types.Config{}
	handlerStrategy := &fakeHandlerStrategy{handlers: map[string]interfaces.ISubscribeHandler{}}
	service := subscriptionservice.NewSubscriptionService(&doneClient{}, cfg)
	service.HandlerStrategy = handlerStrategy

	err := service.Listen(context.Background())
	if err != nil {
		t.Fatalf("Listen deveria terminar sem erro quando não há subscriptions: %v", err)
	}
}
