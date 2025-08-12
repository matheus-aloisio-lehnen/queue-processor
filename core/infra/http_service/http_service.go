package httpservice

import (
	"context"
	"net/http"
	"queue/core/application/subscription/dto"
	"queue/core/domain/types"
	"queue/core/infra/utils/functions"
)

type HttpService struct {
	client functions.IHttpClient
}

func NewHttpService() *HttpService {
	return &HttpService{
		client: &http.Client{},
	}
}

// Opcional para testes – permite substituir o client
func (s *HttpService) WithClient(client functions.IHttpClient) *HttpService {
	s.client = client
	return s
}

func (s *HttpService) SendNotification(data *subscriptiondto.CreateNotificationDto, cfg *types.Config) error {
	url := cfg.URLs.Notification + "/notification"
	headers := functions.CreateBasicAuthHeader(
		cfg.Auth.Username,
		cfg.Auth.Password,
	)
	config := BuildRequestConfig("POST", url, data, nil, nil, headers, true)
	_, err := functions.Send[any](s.client, config)
	return err
}

func BuildRequestConfig(method string, url string, data interface{}, params map[string]string, ctx context.Context, headers map[string]string, showError bool) functions.RequestConfig {
	return functions.RequestConfig{
		Method:    method,
		Url:       url,
		Data:      data,
		Params:    params,
		Ctx:       ctx,
		Headers:   headers,
		ShowError: showError,
	}
}
