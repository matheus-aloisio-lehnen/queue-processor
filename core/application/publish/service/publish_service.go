package publishservice

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"log"
	"queue/core/application/publish/dto"
	"queue/core/domain/interfaces"
)

type PublishService struct {
	Client interfaces.IPubSubClient
}

func NewPublishService(client interfaces.IPubSubClient) (*PublishService, error) {
	return &PublishService{Client: client}, nil
}

func (ps *PublishService) Publish(dto publishdto.InputDto) (bool, error) {
	data, err := json.Marshal(dto)
	if err != nil {
		return false, err
	}
	result := ps.Client.
		Topic(dto.Meta.Topic).
		Publish(context.Background(), &pubsub.Message{Data: data})
	id, err := result.Get(context.Background())
	if err != nil {
		return false, err
	}
	log.Printf("Mensagem publicada com sucesso. ID: %s", id)
	return true, nil
}
