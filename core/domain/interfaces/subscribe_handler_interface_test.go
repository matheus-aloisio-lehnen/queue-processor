package interfaces_test

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/assert"
	"queue/core/domain/interfaces"
)

type SampleSubscription struct {
	Identifier string
}

func (s *SampleSubscription) ID() string {
	return s.Identifier
}
func (s *SampleSubscription) Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error {
	return nil
}

type SimpleHandler struct {
	Handled     bool
	ShouldError error
}

func (sh *SimpleHandler) Handle(sub interfaces.ISubscription) error {
	sh.Handled = true
	return sh.ShouldError
}

func TestISubscribeHandler_Handle_Success(t *testing.T) {
	handler := &SimpleHandler{}
	sub := &SampleSubscription{Identifier: "sub-x"}

	err := handler.Handle(sub)
	assert.True(t, handler.Handled, "Deveria marcar Handled como true")
	assert.NoError(t, err)
}

func TestISubscribeHandler_Handle_Error(t *testing.T) {
	expectedErr := errors.New("falha simulada")
	handler := &SimpleHandler{ShouldError: expectedErr}
	sub := &SampleSubscription{Identifier: "sub-y"}

	err := handler.Handle(sub)
	assert.True(t, handler.Handled, "Deveria marcar Handled como true")
	assert.EqualError(t, err, "falha simulada")
}
