package service

import (
	accountService "github.com/ce-final-project/backend_game_server/account/proto"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/auth/commands"
	"github.com/ce-final-project/backend_game_server/authentication/internal/auth/queries"
	kafkaClient "github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
)

type AuthService struct {
	Commands *commands.AuthCommands
	Queries  *queries.AuthQueries
}

func NewAuthService(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer, asClient accountService.AccountServiceClient) *AuthService {

	registerCommandHandler := commands.NewRegisterAccountHandler(log, cfg, kafkaProducer)

	getAccountByUsernameHandler := queries.NewGetAccountByUsernameHandler(log, cfg, asClient)
	getAccountByIdHandler := queries.NewGetAccountByIdHandler(log, cfg, asClient)

	authCommands := commands.NewAuthCommands(registerCommandHandler)
	authQueries := queries.NewAuthQueries(getAccountByUsernameHandler, getAccountByIdHandler)
	return &AuthService{
		Commands: authCommands,
		Queries:  authQueries,
	}
}
