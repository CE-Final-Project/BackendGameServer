package commands

import (
	"context"
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/dto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type LoginAccountCmdHandler interface {
	Handle(ctx context.Context, command *LoginAccountCommand) (*dto.LoginAccountResponse, error)
}

type loginAccountHandler struct {
	log       logger.Logger
	cfg       *config.Config
	asClient  authService.AuthServiceClient
	accClient authService.AccountServiceClient
}

func NewLoginAccountHandler(log logger.Logger, cfg *config.Config, asClient authService.AuthServiceClient, accClient authService.AccountServiceClient) LoginAccountCmdHandler {
	return &loginAccountHandler{
		log:       log,
		cfg:       cfg,
		asClient:  asClient,
		accClient: accClient,
	}
}

func (r *loginAccountHandler) Handle(ctx context.Context, command *LoginAccountCommand) (*dto.LoginAccountResponse, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "loginAccountHandler.Handle")
	defer span.Finish()

	loginResult, err := r.asClient.Login(ctx, &authService.LoginReq{
		Username: command.LoginDto.Username,
		Password: command.LoginDto.Password,
	})
	if err != nil {
		return nil, err
	}

	result, err := r.accClient.GetAccountByID(ctx, &authService.GetAccountByIdReq{AccountID: loginResult.GetAccountID()})
	if err != nil {
		return nil, err
	}

	account := result.GetAccount()

	return &dto.LoginAccountResponse{
		Account: dto.Account{
			ID:       account.GetAccountID(),
			PlayerID: account.GetPlayerID(),
			Username: account.GetUsername(),
			Email:    account.GetEmail(),
		},
		Token: loginResult.GetToken(),
	}, nil
}
