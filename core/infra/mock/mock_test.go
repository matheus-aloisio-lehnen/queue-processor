package mock_test

import (
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"log"
	"queue/core/infra/mock"
	"testing"
)

func TestMockSubscription_ID(t *testing.T) {
	mockSub := &mock.MockSubscription{IDValue: "sub-1"}
	if mockSub.ID() != "sub-1" {
		t.Errorf("ID esperado 'sub-1', obteve '%s'", mockSub.ID())
	}
}

func TestMockSubscription_Receive(t *testing.T) {
	mockSub := &mock.MockSubscription{IDValue: "sub-1"}
	called := false
	err := mockSub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		called = true
	})
	if err != nil {
		t.Errorf("Receive deveria retornar nil, retornou erro: %v", err)
	}
	log.Println(called)
	// Como a função passada não é chamada (mock faz nada), called segue falso — e está ok.
}

func TestMockSubscriptionIterator_Next(t *testing.T) {
	it := mock.NewMockSubscriptionIterator("a", "b")
	sub, err := it.Next()
	if err != nil {
		t.Fatalf("Esperava sub válido na primeira chamada, erro: %v", err)
	}
	if sub.ID() != "a" {
		t.Errorf("Esperava id 'a', obteve '%s'", sub.ID())
	}

	sub, err = it.Next()
	if err != nil {
		t.Fatalf("Esperava sub válido na segunda chamada, erro: %v", err)
	}
	if sub.ID() != "b" {
		t.Errorf("Esperava id 'b', obteve '%s'", sub.ID())
	}

	sub, err = it.Next()
	if !errors.Is(err, iterator.Done) {
		t.Errorf("Esperava iterator.Done ao terminar, recebeu: %v", err)
	}
	if sub != nil {
		t.Errorf("Sub deveria ser nil após terminar")
	}
}

func TestMockTopic_ID(t *testing.T) {
	topic := &mock.MockTopic{IDValue: "top-123"}
	if topic.ID() != "top-123" {
		t.Errorf("Esperava id 'top-123', obteve '%s'", topic.ID())
	}
}

func TestMockTopic_Publish(t *testing.T) {
	topic := &mock.MockTopic{IDValue: "top-1"}
	result := topic.Publish(context.Background(), &pubsub.Message{Data: []byte("x")})
	if result == nil {
		t.Errorf("Esperava um PublishResult válido")
	}
}

func TestMockPubSubClientAdapter(t *testing.T) {
	client := mock.NewMockPubSubClientAdapter("sub-a", "sub-b")

	// Testa Subscriptions
	it := client.Subscriptions(context.Background())
	sub, err := it.Next()
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if sub.ID() != "sub-a" {
		t.Errorf("Esperava 'sub-a', obteve '%s'", sub.ID())
	}

	// Testa Topic
	topic := client.Topic("top-x")
	if topic.ID() != "top-x" {
		t.Errorf("Esperava id 'top-x', obteve '%s'", topic.ID())
	}

	if err := client.Close(); err != nil {
		t.Errorf("Close deveria retornar nil, retornou: %v", err)
	}
}
