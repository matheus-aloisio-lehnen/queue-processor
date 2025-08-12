package mock

import (
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"queue/core/domain/interfaces"
)

// MockSubscription implementa interfaces.ISubscription
type MockSubscription struct {
	IDValue string
}

func (m *MockSubscription) ID() string { return m.IDValue }

func (m *MockSubscription) Receive(ctx context.Context, fn func(context.Context, *pubsub.Message)) error {
	return nil
}

// MockSubscriptionIterator implementa interfaces.ISubscriptionIterator
type MockSubscriptionIterator struct {
	index int
	ids   []string
}

func NewMockSubscriptionIterator(ids ...string) *MockSubscriptionIterator {
	return &MockSubscriptionIterator{
		index: 0,
		ids:   ids,
	}
}

func (it *MockSubscriptionIterator) Next() (interfaces.ISubscription, error) {
	if it.index < len(it.ids) {
		sub := &MockSubscription{IDValue: it.ids[it.index]}
		it.index++
		return sub, nil
	}
	return nil, iterator.Done
}

type MockPublishResult struct {
	MsgID string
	Err   error
}

func (m *MockPublishResult) Get(ctx context.Context) (string, error) {
	return m.MsgID, m.Err
}

type MockTopic struct {
	IDValue string
}

func (m *MockTopic) ID() string { return m.IDValue }

func (m *MockTopic) Stop() {}

func (m *MockTopic) Publish(ctx context.Context, msg *pubsub.Message) interfaces.IPublishResult {
	if m.IDValue == "fail" {
		return &MockPublishResult{MsgID: "", Err: errors.New("erro simulado")}
	}
	return &MockPublishResult{MsgID: "mocked-message-id", Err: nil}
}

// MockPubSubClientAdapter implementa interfaces.IPubSubClient
type MockPubSubClientAdapter struct {
	SubscriptionIDs []string
}

func NewMockPubSubClientAdapter(subscriptionIDs ...string) interfaces.IPubSubClient {
	return &MockPubSubClientAdapter{SubscriptionIDs: subscriptionIDs}
}

func (c *MockPubSubClientAdapter) Subscriptions(ctx context.Context) interfaces.ISubscriptionIterator {
	ids := c.SubscriptionIDs
	if len(ids) == 0 {
		ids = []string{"test-sub"}
	}
	return NewMockSubscriptionIterator(ids...)
}

func (c *MockPubSubClientAdapter) Topic(id string) interfaces.ITopic {
	return &MockTopic{IDValue: id}
}

func (c *MockPubSubClientAdapter) Close() error {
	return nil
}

type mockSubscriptionStrategy struct{}

func (m *mockSubscriptionStrategy) GetHandler(topic string) interfaces.ISubscribeHandler {
	return nil
}

type MockSubscribeHandler struct {
	CalledWith interfaces.ISubscription
	ShouldFail bool
}

func (m *MockSubscribeHandler) Handle(sub interfaces.ISubscription) error {
	m.CalledWith = sub
	if m.ShouldFail {
		return errors.New("erro simulado")
	}
	return nil
}

type FakeHandlerStrategy struct {
	Handlers map[string]interfaces.ISubscribeHandler
}

func (f *FakeHandlerStrategy) GetHandler(topic string) interfaces.ISubscribeHandler {
	return f.Handlers[topic]
}
