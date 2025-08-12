package subscriptionservice

import (
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"log"
	"queue/core/domain/interfaces"
	"queue/core/domain/strategy"
	"queue/core/domain/types"
	"queue/core/infra/adapter"
)

type SubscriptionService struct {
	Client          interfaces.IPubSubClient
	Cfg             *types.Config
	HandlerStrategy interfaces.HandlerStrategyInterface
}

func NewSubscriptionService(client interfaces.IPubSubClient, cfg *types.Config) *SubscriptionService {
	handlerStrategy := strategy.NewSubscriptionHandlerStrategy(cfg)
	return &SubscriptionService{
		Client:          client,
		Cfg:             cfg,
		HandlerStrategy: handlerStrategy,
	}
}

func (l *SubscriptionService) Listen(ctx context.Context) error {
	subs := l.Client.Subscriptions(ctx)

	for {
		sub, err := subs.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
			log.Printf("Erro ao iterar sobre as subscriptions: %v", err)
			return err
		}
		handler := l.HandlerStrategy.GetHandler(sub.ID())
		if handler == nil {
			log.Printf("Nenhum handler registrado para o tópico %s", sub.ID())
			continue
		}
		log.Printf("Handler encontrado para o tópico %s", sub.ID())
		subAdapter := pubsubadapter.NewSubscriptionAdapter(sub)
		err = handler.Handle(subAdapter)
		if err != nil {
			return errors.New("Erro ao processar assinatura")
		}
	}
	return nil
}
