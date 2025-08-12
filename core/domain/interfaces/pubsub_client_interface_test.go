package interfaces_test

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/assert"
	"queue/core/domain/interfaces"
)

type DemoSubscription struct {
	tag string
}

func (d *DemoSubscription) ID() string {
	return d.tag
}
func (d *DemoSubscription) Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error {
	return nil
}

type ResultStub struct {
	id  string
	err error
}

func (r *ResultStub) Get(ctx context.Context) (string, error) {
	return r.id, r.err
}

type ChannelTopic struct {
	name     string
	failSend bool
}

func (c *ChannelTopic) ID() string { return c.name }
func (c *ChannelTopic) Stop()      {}
func (c *ChannelTopic) Publish(ctx context.Context, msg *pubsub.Message) interfaces.IPublishResult {
	if c.failSend {
		return &ResultStub{id: "", err: errors.New("publish error")}
	}
	return &ResultStub{id: "id-" + c.name, err: nil}
}

type DemoIterator struct {
	index int
	items []interfaces.ISubscription
}

func (it *DemoIterator) Next() (interfaces.ISubscription, error) {
	if it.index < len(it.items) {
		s := it.items[it.index]
		it.index++
		return s, nil
	}
	return nil, errors.New("no more")
}

type DemoPubSubClient struct {
	topics        map[string]*ChannelTopic
	subscriptions []interfaces.ISubscription
	closed        bool
}

func (d *DemoPubSubClient) Subscriptions(ctx context.Context) interfaces.ISubscriptionIterator {
	return &DemoIterator{items: d.subscriptions}
}
func (d *DemoPubSubClient) Topic(id string) interfaces.ITopic {
	t, ok := d.topics[id]
	if ok {
		return t
	}
	return &ChannelTopic{name: id}
}
func (d *DemoPubSubClient) Close() error {
	d.closed = true
	return nil
}

func TestDemoPubSubClient_Flow(t *testing.T) {
	s1 := &DemoSubscription{tag: "A"}
	s2 := &DemoSubscription{tag: "B"}
	topicFoo := &ChannelTopic{name: "foo"}
	client := &DemoPubSubClient{
		topics: map[string]*ChannelTopic{
			"foo": topicFoo,
		},
		subscriptions: []interfaces.ISubscription{s1, s2},
	}

	// Subscriptions iterator
	iter := client.Subscriptions(context.Background())
	sub, err := iter.Next()
	assert.NoError(t, err)
	assert.Equal(t, "A", sub.ID())
	sub2, err := iter.Next()
	assert.NoError(t, err)
	assert.Equal(t, "B", sub2.ID())
	_, err = iter.Next()
	assert.Error(t, err)

	// Topic and publish
	topic := client.Topic("foo")
	assert.Equal(t, "foo", topic.ID())
	res := topic.Publish(context.Background(), &pubsub.Message{Data: []byte("test")})
	id, err := res.Get(context.Background())
	assert.NoError(t, err)
	assert.Contains(t, id, "id-foo")

	// Publish error route
	failTopic := &ChannelTopic{name: "fail", failSend: true}
	failRes := failTopic.Publish(context.Background(), &pubsub.Message{})
	_, err = failRes.Get(context.Background())
	assert.EqualError(t, err, "publish error")

	// Close
	assert.NoError(t, client.Close())
	assert.True(t, client.closed)
}
