package config

import (
	"github.com/gin-contrib/cors"
	"os"
	"queue/core/domain/types"
	"strconv"
	"time"
)

func LoadConfig() *types.Config {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	return &types.Config{
		Environment: os.Getenv("ENVIRONMENT"),
		Port:        port,
		CorsConfig: cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"*"},
			AllowCredentials: false,
			ExposeHeaders:    []string{"Content-Length"},
			MaxAge:           12 * time.Hour,
		},
		Auth: types.BasicAuthConfig{
			Username: os.Getenv("BASIC_AUTH_USERNAME"),
			Password: os.Getenv("BASIC_AUTH_PASSWORD"),
		},
		Google: types.GoogleConfig{
			ProjectID:   os.Getenv("PROJECT_ID"),
			Credentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		},
		URLs: types.URLsConfig{
			Frontend:     os.Getenv("FRONTEND_URL"),
			API:          os.Getenv("API_URL"),
			Notification: os.Getenv("NOTIFICATION_URL"),
			Storage:      os.Getenv("STORAGE_URL"),
			Queue:        os.Getenv("QUEUE_URL"),
		},
	}
}
