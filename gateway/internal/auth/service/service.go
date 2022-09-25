package service

import (
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/auth/commands"
	"github.com/ce-final-project/backend_game_server/gateway/internal/auth/queries"
	kafkaClient "github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
)

type AuthService struct {
	Commands *commands.AuthCommands
	Queries  *queries.AuthQueries
}

func NewAuthService(
	log logger.Logger,
	cfg *config.Config,
	kafkaProducer kafkaClient.Producer,
	authService authService.AuthServiceClient,
) *AuthService {

	registerCommandHandler := commands.NewRegisterAccountHandler(log, cfg, authService)
	loginCommandHandler := commands.NewLoginAccountHandler(log, cfg, authService)

	verifyTokenQueryHandler := queries.NewVerifyTokenHandler(log, cfg, authService)

	authCommands := commands.NewAuthCommands(registerCommandHandler, loginCommandHandler)
	authQueries := queries.NewAccountQueries(verifyTokenQueryHandler)

	return &AuthService{
		Commands: authCommands,
		Queries:  authQueries,
	}
}
