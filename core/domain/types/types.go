package types

import "github.com/gin-contrib/cors"

type BasicAuthConfig struct {
	Username string
	Password string
}

type GoogleConfig struct {
	ProjectID   string
	Credentials string
}

type URLsConfig struct {
	Frontend     string
	API          string
	Notification string
	Storage      string
	Queue        string
}

type Config struct {
	Environment string
	Port        int
	CorsConfig  cors.Config
	Auth        BasicAuthConfig
	Google      GoogleConfig
	URLs        URLsConfig
}
