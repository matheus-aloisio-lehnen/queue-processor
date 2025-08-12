package publishservice_test

import (
	"context"
	"encoding/json"
	"errors"
	"queue/core/application/publish/dto"
	"queue/core/application/publish/service"
	"queue/core/domain/interfaces"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para IPubSubClient
type MockPubSubClient struct {
	mock.Mock
}

func (m *MockPubSubClient) Topic(name string) interfaces.ITopic {
	args := m.Called(name)
	return args.Get(0).(interfaces.ITopic)
}

func (m *MockPubSubClient) Subscriptions(ctx context.Context) interfaces.ISubscriptionIterator {
	args := m.Called(ctx)
	return args.Get(0).(interfaces.ISubscriptionIterator)
}

func (m *MockPubSubClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Mock para ITopic
type MockPubSubTopic struct {
	mock.Mock
}

func (m *MockPubSubTopic) ID() string { return "" }
func (m *MockPubSubTopic) Stop()      {}
func (m *MockPubSubTopic) Publish(ctx context.Context, msg *pubsub.Message) interfaces.IPublishResult {
	args := m.Called(ctx, msg)
	return args.Get(0).(interfaces.IPublishResult)
}

type MockPublishResult struct {
	mock.Mock
}

func (m *MockPublishResult) Get(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func TestPublish_Success(t *testing.T) {
	clientMock := new(MockPubSubClient)
	topicMock := new(MockPubSubTopic)
	publishResult := new(MockPublishResult)

	dtoInput := publishdto.InputDto{
		Meta: publishdto.MetaDto{Topic: "my-topic"},
		Data: map[string]string{"key": "value"},
	}
	data, _ := json.Marshal(dtoInput)

	// Configura mocks
	topicMock.On("Publish", mock.Anything, mock.MatchedBy(func(msg *pubsub.Message) bool {
		return string(msg.Data) == string(data)
	})).Return(publishResult)
	publishResult.On("Get", mock.Anything).Return("msg-id", nil)
	clientMock.On("Topic", "my-topic").Return(topicMock)

	svc := &publishservice.PublishService{Client: clientMock}
	ok, err := svc.Publish(dtoInput)
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestPublish_JSONMarshalError(t *testing.T) {
	clientMock := new(MockPubSubClient)
	dtoInput := publishdto.InputDto{
		Meta: publishdto.MetaDto{Topic: "any"},
		Data: func() {},
	}
	svc := &publishservice.PublishService{Client: clientMock}
	ok, err := svc.Publish(dtoInput)
	assert.False(t, ok)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json:")
}

func TestPublish_PubSubError(t *testing.T) {
	clientMock := new(MockPubSubClient)
	topicMock := new(MockPubSubTopic)
	publishResult := new(MockPublishResult)

	dtoInput := publishdto.InputDto{
		Meta: publishdto.MetaDto{Topic: "erro-topic"},
		Data: map[string]string{"k": "v"},
	}
	data, _ := json.Marshal(dtoInput)

	// Configura mocks
	topicMock.On("Publish", mock.Anything, mock.MatchedBy(func(msg *pubsub.Message) bool {
		return string(msg.Data) == string(data)
	})).Return(publishResult)
	publishResult.On("Get", mock.Anything).Return("", errors.New("erro pub sub"))
	clientMock.On("Topic", "erro-topic").Return(topicMock)

	svc := &publishservice.PublishService{Client: clientMock}
	ok, err := svc.Publish(dtoInput)
	assert.False(t, ok)
	assert.EqualError(t, err, "erro pub sub")
}
