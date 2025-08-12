package interfaces

import (
	"cloud.google.com/go/pubsub"
	"context"
)

type IPubSubClient interface {
	Subscriptions(ctx context.Context) ISubscriptionIterator
	Topic(id string) ITopic
	Close() error
}

type ISubscriptionIterator interface {
	Next() (ISubscription, error)
}

type ITopic interface {
	ID() string
	Stop()
	Publish(ctx context.Context, msg *pubsub.Message) IPublishResult
}

type ISubscription interface {
	ID() string
	Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error
}

type IPublishResult interface {
	Get(ctx context.Context) (string, error)
}
