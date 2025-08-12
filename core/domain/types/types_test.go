package types_test

import (
	"github.com/stretchr/testify/assert"
	"queue/core/domain/types"
	"testing"
)

func TestBasicAuthConfig(t *testing.T) {
	auth := types.BasicAuthConfig{
		Username: "admin",
		Password: "secret",
	}

	assert.Equal(t, "admin", auth.Username)
	assert.Equal(t, "secret", auth.Password)
}

func TestGoogleConfig(t *testing.T) {
	googleConfig := types.GoogleConfig{
		ProjectID:   "project-id",
		Credentials: "credentials-file",
	}
	assert.Equal(t, "project-id", googleConfig.ProjectID)
	assert.Equal(t, "credentials-file", googleConfig.Credentials)
}

func TestURLsConfig(t *testing.T) {
	urlsConfig := types.URLsConfig{
		Frontend:     "http://localhost",
		API:          "http://localhost/api",
		Notification: "http://localhost/notifications",
		Storage:      "http://localhost/storage",
		Queue:        "http://localhost/queue",
	}

	assert.Equal(t, "http://localhost", urlsConfig.Frontend)
	assert.Equal(t, "http://localhost/api", urlsConfig.API)
	assert.Equal(t, "http://localhost/notifications", urlsConfig.Notification)
	assert.Equal(t, "http://localhost/storage", urlsConfig.Storage)
	assert.Equal(t, "http://localhost/queue", urlsConfig.Queue)
}

func TestConfig(t *testing.T) {
	cfg := types.Config{
		Environment: "development",
		Port:        3003,
		Auth: types.BasicAuthConfig{
			Username: "admin",
			Password: "password",
		},
		Google: types.GoogleConfig{
			ProjectID:   "project-id",
			Credentials: "credentials-file",
		},
		URLs: types.URLsConfig{
			Frontend:     "http://localhost",
			API:          "http://localhost/api",
			Notification: "http://localhost/notifications",
			Storage:      "http://localhost/storage",
			Queue:        "http://localhost/queue",
		},
	}
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, 3003, cfg.Port)
	assert.Equal(t, "admin", cfg.Auth.Username)
	assert.Equal(t, "password", cfg.Auth.Password)
	assert.Equal(t, "project-id", cfg.Google.ProjectID)
	assert.Equal(t, "credentials-file", cfg.Google.Credentials)
	assert.Equal(t, "http://localhost", cfg.URLs.Frontend)
	assert.Equal(t, "http://localhost/api", cfg.URLs.API)
	assert.Equal(t, "http://localhost/notifications", cfg.URLs.Notification)
	assert.Equal(t, "http://localhost/storage", cfg.URLs.Storage)
	assert.Equal(t, "http://localhost/queue", cfg.URLs.Queue)
}
