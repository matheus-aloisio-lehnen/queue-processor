package pubsubadapter_test

import (
	"context"
	"errors"
	"queue/core/infra/adapter"
	"testing"

	"cloud.google.com/go/pubsub"
	"queue/core/domain/interfaces"
)

// --- Mocks ---

type mockSubscription struct {
	id          string
	receiveFunc func(ctx context.Context, f func(context.Context, *pubsub.Message)) error
}

func (m *mockSubscription) ID() string {
	return m.id
}
func (m *mockSubscription) Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error {
	if m.receiveFunc != nil {
		return m.receiveFunc(ctx, f)
	}
	return nil
}

type mockSubscriptionIterator struct {
	subs  []interfaces.ISubscription
	index int
}

func (m *mockSubscriptionIterator) Next() (interfaces.ISubscription, error) {
	if m.index >= len(m.subs) {
		return nil, errors.New("done")
	}
	sub := m.subs[m.index]
	m.index++
	return sub, nil
}

// --- Test SubscriptionAdapter ---

func TestSubscriptionAdapter_ID(t *testing.T) {
	mockSub := &mockSubscription{id: "test-sub"}
	adapter := pubsubadapter.NewSubscriptionAdapter(mockSub)

	if adapter.ID() != "test-sub" {
		t.Errorf("expected sub ID to be 'test-sub', got '%s'", adapter.ID())
	}
}

func TestSubscriptionAdapter_Receive(t *testing.T) {
	called := false
	f := func(ctx context.Context, msg *pubsub.Message) {
		called = true
	}

	mockSub := &mockSubscription{
		id: "x",
		receiveFunc: func(ctx context.Context, recvFunc func(context.Context, *pubsub.Message)) error {
			recvFunc(ctx, &pubsub.Message{})
			return nil
		},
	}
	adapter := pubsubadapter.NewSubscriptionAdapter(mockSub)
	err := adapter.Receive(context.Background(), f)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !called {
		t.Error("expected receive function to be called")
	}
}
