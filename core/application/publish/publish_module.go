package publishmodule

import (
	"errors"
	"github.com/gin-gonic/gin"
	"queue/core/application/publish/controller"
	"queue/core/application/publish/dto"
	"queue/core/application/publish/service"
	"queue/core/domain/interfaces"
	"queue/core/infra/middleware"
)

type PublishModule struct {
	Controller *publishcontroller.PublishController
}

func NewPublishModule(client interfaces.IPubSubClient) (*PublishModule, error) {
	publishService, err := publishservice.NewPublishService(client)
	if err != nil {
		return nil, errors.New("Erro ao criar o serviço de publicação")
	}
	controller := publishcontroller.NewPublishController(*publishService)
	return &PublishModule{
		Controller: controller,
	}, nil
}

func (m *PublishModule) RegisterRoutes(router *gin.Engine) {
	publishGroup := router.Group("/publish")
	{
		publishGroup.POST("", classtransformer.UseClassTransformerMiddleware(&publishdto.InputDto{}), m.Controller.Publish)
	}
}
