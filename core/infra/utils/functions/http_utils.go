package functions

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type IHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestConfig struct {
	Method    string
	Url       string
	Data      any
	Params    map[string]string
	Headers   map[string]string
	Ctx       context.Context
	ShowError bool
}

type HttpException struct {
	Message string
	Status  int
}

func (e *HttpException) Error() string {
	return fmt.Sprintf("status %d: %s", e.Status, e.Message)
}

func Send[T any](client IHttpClient, cfg RequestConfig) (*T, error) {
	req, err := buildRequest(cfg)
	if err != nil {
		return returnError[T](cfg.ShowError, "Erro ao criar requisição", err, http.StatusInternalServerError)
	}

	resp, err := client.Do(req)
	if err != nil {
		return returnError[T](cfg.ShowError, "Erro ao executar requisição", err, http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	return handleResponse[T](resp, cfg.ShowError)
}

// ---------- Helpers internos ----------

func buildRequest(cfg RequestConfig) (*http.Request, error) {
	u, err := buildURLWithParams(cfg.Url, cfg.Params)
	if err != nil {
		return nil, err
	}
	var body io.Reader
	if cfg.Data != nil && allowsBody(cfg.Method) {
		b, err := json.Marshal(cfg.Data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(b)
	}
	ctx := cfg.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	req, err := http.NewRequestWithContext(ctx, cfg.Method, u, body)
	if err != nil {
		return nil, err
	}
	// Headers
	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}
	if cfg.Data != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func buildURLWithParams(rawURL string, params map[string]string) (string, error) {
	if len(params) == 0 {
		return rawURL, nil
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func allowsBody(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return true
	}
	return false
}

func handleResponse[T any](resp *http.Response, showError bool) (*T, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return returnError[T](showError, "Erro ao ler resposta", err, http.StatusInternalServerError)
	}
	if resp.StatusCode >= 400 {
		if !showError {
			return nil, nil
		}
		return nil, parseHttpError(body, resp.StatusCode)
	}
	var result T
	if len(body) > 0 {
		if err := json.Unmarshal(body, &result); err != nil {
			return returnError[T](showError, "Erro ao decodificar JSON", err, http.StatusInternalServerError)
		}
	}
	return &result, nil
}

func parseHttpError(body []byte, status int) error {
	var fromAPI struct{ Message string }
	msg := ""
	if len(body) > 0 && json.Unmarshal(body, &fromAPI) == nil && fromAPI.Message != "" {
		msg = fromAPI.Message
	}
	if msg == "" {
		msg = "Erro desconhecido"
	}
	return &HttpException{
		Message: msg,
		Status:  status,
	}
}

func returnError[T any](showError bool, msg string, err error, status int) (*T, error) {
	if !showError {
		return nil, nil
	}
	if err != nil && err.Error() != "" {
		msg = fmt.Sprintf("%s: %s", msg, err.Error())
	}
	return nil, &HttpException{Message: msg, Status: status}
}

func CreateBasicAuthHeader(username, password string) map[string]string {
	credentials := fmt.Sprintf("%s:%s", username, password)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return map[string]string{
		"Authorization": "Basic " + encoded,
	}
}
