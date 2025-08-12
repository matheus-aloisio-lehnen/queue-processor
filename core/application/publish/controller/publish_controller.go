package publishcontroller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"queue/core/application/publish/dto"
	"queue/core/application/publish/service"
	"queue/core/domain/response"
)

type PublishController struct {
	service publishservice.PublishService
}

func NewPublishController(service publishservice.PublishService) *PublishController {
	return &PublishController{
		service: service,
	}
}

func (ctrl *PublishController) Publish(c *gin.Context) {
	dto, exists := c.Get("dto")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Formato de mensagem inválido", nil)
		return
	}
	messageDto, ok := dto.(*publishdto.InputDto)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Formato de mensagem inválido", nil)
		return
	}
	_, err := ctrl.service.Publish(*messageDto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), err)
		return
	}
	response.Success(c, true, 200, "Mensagem publicada com sucesso")
}
