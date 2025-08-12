package interfaces

type HandlerStrategyInterface interface {
	GetHandler(topic string) ISubscribeHandler
}
