package commands

import (
	"context"
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/dto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	uuid "github.com/satori/go.uuid"
)

type RegisterAccountCmdHandler interface {
	Handle(ctx context.Context, command *RegisterAccountCommand) (*dto.RegisterAccountResponse, error)
}

type registerAccountHandler struct {
	log      logger.Logger
	cfg      *config.Config
	asClient authService.AuthServiceClient
}

func NewRegisterAccountHandler(log logger.Logger, cfg *config.Config, asClient authService.AuthServiceClient) *registerAccountHandler {
	return &registerAccountHandler{
		log:      log,
		cfg:      cfg,
		asClient: asClient,
	}
}

func (r *registerAccountHandler) Handle(ctx context.Context, command *RegisterAccountCommand) (*dto.RegisterAccountResponse, error) {

	regResult, err := r.asClient.Register(ctx, &authService.RegisterReq{
		Username: command.RegisterDto.Username,
		Email:    command.RegisterDto.Email,
		Password: command.RegisterDto.Password,
	})
	if err != nil {
		return nil, err
	}

	accountUUID, err := uuid.FromString(regResult.GetAccountID())
	if err != nil {
		return nil, err
	}

	return &dto.RegisterAccountResponse{
		AccountID: accountUUID,
		PlayerID:  regResult.GetPlayerID(),
		Token:     regResult.GetToken(),
	}, nil
}