package main

import (
	"github.com/ce-final-project/backend_game_server/account/config"
	"github.com/ce-final-project/backend_game_server/account/internal/server"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	cfg, err := config.InitConfig("")
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("AccountService")

	s := server.NewServer(appLogger, cfg)
	appLogger.Fatal(s.Run())
}
