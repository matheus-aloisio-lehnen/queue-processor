package servers

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"queue/core/application/auth"
	"queue/core/application/health_check"
	"queue/core/application/publish"
	"queue/core/application/subscription"
	"queue/core/domain/interfaces"
	"queue/core/domain/types"
	"queue/core/infra/adapter"
	"queue/core/infra/config"
	"queue/core/infra/exceptions"
)

func Run() {
	LoadEnv()
	cfg := config.LoadConfig()

	pubsubClient := NewPubSubClient(cfg)
	adapter := pubsubadapter.NewPubSubClientAdapter(pubsubClient)

	router := SetupRouter(cfg, adapter)
	if err := StartServer(router, cfg.Port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Nenhum arquivo .env encontrado. Usando variáveis de ambiente do sistema: %v", err)
	}
}

func SetupRouter(cfg *types.Config, pubsubClient interfaces.IPubSubClient) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cfg.CorsConfig), exceptions.AllExceptionFilter())
	r.Use(auth.BasicAuthMiddleware(cfg.Auth.Username, cfg.Auth.Password))
	RegisterModules(r, cfg, pubsubClient)
	return r
}

func RegisterModules(r *gin.Engine, cfg *types.Config, pubsubClient interfaces.IPubSubClient) {
	healthCheckModule := healthcheckmodule.NewHealthCheckModule()
	healthCheckModule.RegisterRoutes(r)
	publishModule, err := publishmodule.NewPublishModule(pubsubClient)
	if err != nil {
		log.Fatalf("Erro ao criar publish module: %v", err)
	}
	publishModule.RegisterRoutes(r)
	subscriptionService := subscriptionservice.NewSubscriptionService(pubsubClient, cfg)
	go func() {
		if err := subscriptionService.Listen(context.Background()); err != nil {
			log.Fatalf("Erro no listener: %v", err)
		}
	}()
}

func NewPubSubClient(cfg *types.Config) *pubsub.Client {
	client, err := pubsub.NewClient(context.Background(), cfg.Google.ProjectID)
	if err != nil {
		log.Fatalf("Erro ao criar o pubsub: %v", err)
	}
	return client
}

func StartServer(router *gin.Engine, port int) error {
	if port == 0 {
		port = 3003
	}
	log.Printf("Server running on port %d", port)
	return router.Run(fmt.Sprintf(":%d", port))
}
