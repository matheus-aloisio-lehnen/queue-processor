package pubsubadapter

import (
	"cloud.google.com/go/pubsub"
	"context"
	"queue/core/domain/interfaces"
)

type PubSubClientAdapter struct {
	client *pubsub.Client
}

func NewPubSubClientAdapter(client *pubsub.Client) *PubSubClientAdapter {
	return &PubSubClientAdapter{client: client}
}

func (a *PubSubClientAdapter) Subscriptions(ctx context.Context) interfaces.ISubscriptionIterator {
	return &subscriptionIteratorAdapter{a.client.Subscriptions(ctx)}
}

func (a *PubSubClientAdapter) Topic(id string) interfaces.ITopic {
	return &topicAdapter{a.client.Topic(id)}
}

func (a *PubSubClientAdapter) Close() error {
	return a.client.Close()
}

/* ------------------------------------- ITERATORS--------------------------------------------------------------*/

type subscriptionIteratorAdapter struct {
	iter *pubsub.SubscriptionIterator
}

func (s *subscriptionIteratorAdapter) Next() (interfaces.ISubscription, error) {
	return s.iter.Next()
}

type topicAdapter struct {
	topic *pubsub.Topic
}

func (t *topicAdapter) ID() string {
	return t.topic.ID()
}

func (t *topicAdapter) Stop() {
	t.topic.Stop()
}

func (t *topicAdapter) Publish(ctx context.Context, msg *pubsub.Message) interfaces.IPublishResult {
	return t.topic.Publish(ctx, msg)
}

/*--------------------------------------- RECEIVE --------------------------------------------------------------*/

type SubscriptionAdapter struct {
	sub interfaces.ISubscription
}

func NewSubscriptionAdapter(sub interfaces.ISubscription) *SubscriptionAdapter {
	return &SubscriptionAdapter{sub: sub}
}

func (s *SubscriptionAdapter) ID() string {
	return s.sub.ID()
}

func (s *SubscriptionAdapter) Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error {
	return s.sub.Receive(ctx, f)
}

/*---------------------------------				RESULT ADAPTER			-------------------------------------*/

type PublishResultAdapter struct {
	result *pubsub.PublishResult
}

func (a *PublishResultAdapter) Get(ctx context.Context) (string, error) {
	return a.result.Get(ctx)
}
