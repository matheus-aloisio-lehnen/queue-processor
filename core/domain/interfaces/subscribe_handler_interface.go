package interfaces

type ISubscribeHandler interface {
	Handle(sub ISubscription) error
}
