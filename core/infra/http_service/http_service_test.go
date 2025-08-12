package httpservice_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"queue/core/application/subscription/dto"
	"queue/core/domain/types"
	"queue/core/infra/http_service"
	"testing"
)

type MockHttpClient struct {
	resp        *http.Response
	err         error
	lastRequest *http.Request
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	m.lastRequest = req
	return m.resp, m.err
}

func TestSendNotification_Success(t *testing.T) {
	mockClient := &MockHttpClient{
		resp: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
		},
		err: nil,
	}

	service := httpservice.NewHttpService().WithClient(mockClient)
	cfg := &types.Config{
		Auth: types.BasicAuthConfig{
			Username: "admin",
			Password: "123456",
		},
		URLs: types.URLsConfig{
			Notification: "http://localhost:1234",
		},
	}
	nData := &subscriptiondto.CreateNotificationDto{
		UserID:    81,
		UserName:  "Tester",
		Channel:   "EMAIL",
		Recipient: "teste@ddd.com",
		Subject:   "S",
		Payload:   map[string]string{"k": "v"},
	}

	err := service.SendNotification(nData, cfg)
	if err != nil {
		t.Errorf("esperava nil, recebeu erro: %v", err)
	}
	// Confirma que foi chamado POST correto
	if mockClient.lastRequest == nil || mockClient.lastRequest.Method != "POST" {
		t.Errorf("esperava método POST")
	}
	if mockClient.lastRequest.URL.String() != "http://localhost:1234/notification" {
		t.Errorf("url inesperada: %s", mockClient.lastRequest.URL.String())
	}
}

func TestSendNotification_ErroNetwork(t *testing.T) {
	mockClient := &MockHttpClient{
		resp: &http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewBufferString(`{"error":"fail"}`)),
		},
		err: nil,
	}

	service := httpservice.NewHttpService().WithClient(mockClient)
	cfg := &types.Config{
		Auth: types.BasicAuthConfig{
			Username: "admin",
			Password: "123456",
		},
		URLs: types.URLsConfig{
			Notification: "http://localhost:8888",
		},
	}
	nData := &subscriptiondto.CreateNotificationDto{
		UserID:    90,
		UserName:  "Tester",
		Channel:   "EMAIL",
		Recipient: "userx@ddd.com",
		Subject:   "Erro",
		Payload:   map[string]string{"k": "v"},
	}

	err := service.SendNotification(nData, cfg)
	if err == nil {
		t.Error("esperava erro de status HTTP, recebeu nil")
	}

}

func TestBuildRequestConfig(t *testing.T) {
	headers := map[string]string{"Authorization": "Basic x"}
	params := map[string]string{"p": "v"}
	data := "body"
	ctx := context.Background()
	cfg := httpservice.BuildRequestConfig("GET", "http://u", data, params, ctx, headers, true)
	if cfg.Method != "GET" || cfg.Url != "http://u" || cfg.Data != data || cfg.Headers["Authorization"] != "Basic x" || cfg.Params["p"] != "v" || cfg.Ctx != ctx {
		t.Error("buildRequestConfig não montou corretamente a struct")
	}
}
