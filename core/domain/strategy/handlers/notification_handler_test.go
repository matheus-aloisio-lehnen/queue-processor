package handlers_test

import (
	"cloud.google.com/go/pubsub"
	"context"
	"queue/core/domain/strategy/handlers"
	"testing"

	"queue/core/domain/types"
)

type fakeSubscription struct {
	callback func(context.Context, *pubsub.Message)
}

func (f *fakeSubscription) ID() string { return "fake-sub" }
func (f *fakeSubscription) Receive(ctx context.Context, fn func(context.Context, *pubsub.Message)) error {
	f.callback = fn
	return nil
}

func TestNotificationHandler_Handle_ValidMessage(t *testing.T) {
	sub := &fakeSubscription{}
	handler := handlers.NewNotificationHandler(&types.Config{})
	err := handler.Handle(sub)
	if err != nil {
		t.Fatalf("Handler.Handle error: %v", err)
	}
	msg := &pubsub.Message{
		Data: []byte(`{"data":{"userId":1,"channel":"email","recipient":"foo@bar.com","payload":{"sub":"ok"}}}`),
	}
	sub.callback(context.Background(), msg)
}

func TestNotificationHandler_Handle_InvalidMessage(t *testing.T) {
	sub := &fakeSubscription{}
	handler := handlers.NewNotificationHandler(&types.Config{})
	_ = handler.Handle(sub)
	msg := &pubsub.Message{
		Data: []byte(`{not-json}`),
	}
	sub.callback(context.Background(), msg)
}
