package config_test

import (
	"os"
	"queue/core/infra/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("PORT", "3003")
	os.Setenv("BASIC_AUTH_USERNAME", "admin")
	os.Setenv("BASIC_AUTH_PASSWORD", "senha")
	os.Setenv("PROJECT_ID", "projeto-123")
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/caminho/credencial.json")
	os.Setenv("FRONTEND_URL", "http://frontend")
	os.Setenv("API_URL", "http://api")
	os.Setenv("NOTIFICATION_URL", "http://notification")
	os.Setenv("STORAGE_URL", "http://storage")
	os.Setenv("QUEUE_URL", "http://queue")
	cfg := config.LoadConfig()
	assert.Equal(t, "test", cfg.Environment)
	assert.Equal(t, 3003, cfg.Port)
	assert.Equal(t, "admin", cfg.Auth.Username)
	assert.Equal(t, "senha", cfg.Auth.Password)
	assert.Equal(t, "projeto-123", cfg.Google.ProjectID)
	assert.Equal(t, "/caminho/credencial.json", cfg.Google.Credentials)
	assert.Equal(t, "http://frontend", cfg.URLs.Frontend)
	assert.Equal(t, "http://api", cfg.URLs.API)
	assert.Equal(t, "http://notification", cfg.URLs.Notification)
	assert.Equal(t, "http://storage", cfg.URLs.Storage)
	assert.Equal(t, "http://queue", cfg.URLs.Queue)
}
