package main

import (
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/server"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"log"
)

// @title           			API Gateway Game Server
// @version         			1.0
// @description     			API Gateway microservices.
// @BasePath  					/api/v1
// @contact.name 				Poomipat Chuealue
// @contact.url 				https://github.com/Poomipat-Ch
// @contact.email 				poomipat002@gmail.com
// @securityDefinitions.apikey 	BearerAuth
// @in header
// @name X-AUTH-TOKEN
func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("ApiGateway")

	s := server.NewServer(appLogger, cfg)
	appLogger.Fatal(s.Run())
}
