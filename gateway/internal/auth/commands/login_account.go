package commands

import (
	"context"
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/dto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
)

type LoginAccountCmdHandler interface {
	Handle(ctx context.Context, command *LoginAccountCommand) (*dto.LoginAccountResponse, error)
}

type loginAccountHandler struct {
	log      logger.Logger
	cfg      *config.Config
	asClient authService.AuthServiceClient
}

func NewLoginAccountHandler(log logger.Logger, cfg *config.Config, asClient authService.AuthServiceClient) LoginAccountCmdHandler {
	return &loginAccountHandler{
		log:      log,
		cfg:      cfg,
		asClient: asClient,
	}
}

func (r *loginAccountHandler) Handle(ctx context.Context, command *LoginAccountCommand) (*dto.LoginAccountResponse, error) {

	loginResult, err := r.asClient.Login(ctx, &authService.LoginReq{
		Username: command.LoginDto.Username,
		Password: command.LoginDto.Password,
	})
	if err != nil {
		return nil, err
	}

	return &dto.LoginAccountResponse{
		AccountID: loginResult.GetAccountID(),
		PlayerID:  loginResult.GetPlayerID(),
		Token:     loginResult.GetToken(),
	}, nil
}
