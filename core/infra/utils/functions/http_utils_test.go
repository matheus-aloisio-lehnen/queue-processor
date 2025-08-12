package functions_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"queue/core/infra/utils/functions"
	"reflect"
	"testing"
)

type mockClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

type respStruct struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func TestSend_Success(t *testing.T) {
	want := respStruct{Name: "go", ID: 22}
	respBody := `{"name":"go", "id":22}`

	client := &mockClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != "POST" {
				t.Fatalf("unexpected method: %s", req.Method)
			}
			if req.URL.Query().Get("q") != "abc" {
				t.Fatalf("query params missing: %v", req.URL.String())
			}
			if req.Header.Get("X-Test") != "true" {
				t.Fatalf("header missing")
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(respBody)),
			}, nil
		},
	}
	cfg := functions.RequestConfig{
		Method:    "POST",
		Url:       "http://example.com?foo=1",
		Params:    map[string]string{"q": "abc"},
		Headers:   map[string]string{"X-Test": "true"},
		Data:      map[string]string{"x": "y"},
		ShowError: true,
	}
	res, err := functions.Send[respStruct](client, cfg)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !reflect.DeepEqual(*res, want) {
		t.Errorf("want=%v got=%v", want, *res)
	}
}

func TestSend_Success_NoBody(t *testing.T) {
	client := &mockClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBuffer(nil)),
			}, nil
		},
	}
	type empty struct{}
	cfg := functions.RequestConfig{
		Method:    "GET",
		Url:       "http://example.com",
		ShowError: true,
	}
	res, err := functions.Send[empty](client, cfg)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res == nil {
		t.Errorf("result should not be nil")
	}
}

func TestSend_buildRequestError(t *testing.T) {
	cfg := functions.RequestConfig{
		Method:    "POST",
		Url:       "::::://invalid_url",
		ShowError: true,
	}
	client := &mockClient{}
	res, err := functions.Send[respStruct](client, cfg)
	if res != nil {
		t.Error("res should be nil")
	}
	var httpErr *functions.HttpException
	if err == nil || !errors.As(err, &httpErr) {
		t.Error("err should be HttpException")
	}
}

func TestSend_DoError(t *testing.T) {
	cfg := functions.RequestConfig{
		Method:    "GET",
		Url:       "http://example.com",
		ShowError: true,
	}
	client := &mockClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("net error")
		},
	}
	res, err := functions.Send[respStruct](client, cfg)
	if res != nil {
		t.Error("res should be nil")
	}
	var httpErr *functions.HttpException
	if err == nil || !errors.As(err, &httpErr) || httpErr.Message == "" {
		t.Error("should wrap http error")
	}
}

func TestSend_HTTP4XX_WithShowError(t *testing.T) {
	respBody := `{"message":"not found"}`
	client := &mockClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewBufferString(respBody)),
			}, nil
		},
	}
	cfg := functions.RequestConfig{
		Method:    "GET",
		Url:       "http://example.com",
		ShowError: true,
	}
	res, err := functions.Send[respStruct](client, cfg)
	if res != nil {
		t.Error("res should be nil")
	}
	var httpErr *functions.HttpException
	if err == nil || !errors.As(err, &httpErr) {
		t.Error("expected HttpException")
	}
	if httpErr.Message != "not found" {
		t.Errorf("want 'not found', got %v", httpErr.Message)
	}
	if httpErr.Status != 404 {
		t.Error("status should be 404")
	}
}

func TestSend_HTTP400_WithoutShowError(t *testing.T) {
	client := &mockClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message":"fail"}`)),
			}, nil
		},
	}
	cfg := functions.RequestConfig{
		Method:    "GET",
		Url:       "http://example.com",
		ShowError: false,
	}
	res, err := functions.Send[respStruct](client, cfg)
	if res != nil {
		t.Error("res should be nil")
	}
	if err != nil {
		t.Errorf("err should be nil, got %#v", err)
	}
}

func TestSend_ReadBodyError(t *testing.T) {
	rcErr := io.ErrUnexpectedEOF
	client := &mockClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			r := io.NopCloser(badReader{readErr: rcErr})
			return &http.Response{StatusCode: 200, Body: r}, nil
		},
	}
	cfg := functions.RequestConfig{
		Method:    "GET",
		Url:       "http://example.com",
		ShowError: true,
	}
	res, err := functions.Send[respStruct](client, cfg)
	if res != nil {
		t.Error("res should be nil")
	}
	var httpErr *functions.HttpException
	if err == nil || !errors.As(err, &httpErr) {
		t.Error("expected HttpException")
	}
}

func TestSend_JSONDecodeError(t *testing.T) {
	badJSON := "{invalid json"
	client := &mockClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(badJSON)),
			}, nil
		},
	}
	cfg := functions.RequestConfig{
		Method:    "GET",
		Url:       "http://example.com",
		ShowError: true,
	}
	res, err := functions.Send[respStruct](client, cfg)
	if res != nil {
		t.Error("res should be nil")
	}
	var httpErr *functions.HttpException
	if err == nil || !errors.As(err, &httpErr) {
		t.Error("expected HttpException")
	}
}

// Simula um erro na leitura do body
type badReader struct {
	readErr error
}

func (b badReader) Read([]byte) (int, error) { return 0, b.readErr }
func (b badReader) Close() error             { return nil }

func TestSend_ContentTypeHeader(t *testing.T) {
	var gotContentType string
	client := &mockClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			gotContentType = req.Header.Get("Content-Type")
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBuffer(nil)),
			}, nil
		},
	}
	cfg := functions.RequestConfig{
		Method:  "POST",
		Url:     "http://example.com",
		Data:    map[string]string{"foo": "bar"},
		Headers: make(map[string]string), // força setar Content-Type
	}
	_, _ = functions.Send[respStruct](client, cfg)
	if gotContentType != "application/json" {
		t.Errorf("Content-Type should be set to application/json, got %v", gotContentType)
	}
}
