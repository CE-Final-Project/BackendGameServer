package commands

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/config"
	"github.com/ce-final-project/backend_game_server/account/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/account/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"github.com/opentracing/opentracing-go"
)

type CreateAccountCmdHandler interface {
	Handle(ctx context.Context, command *CreateAccountCommand) error
}

type createAccountHandler struct {
	log          logger.Logger
	cfg          *config.Config
	postgresRepo repository.Repository
	redisRepo    repository.CacheRepository
}

func NewCreateAccountHandler(log logger.Logger, cfg *config.Config, mongoRepo repository.Repository, redisRepo repository.CacheRepository) CreateAccountCmdHandler {
	return &createAccountHandler{log: log, cfg: cfg, postgresRepo: mongoRepo, redisRepo: redisRepo}
}

func (c *createAccountHandler) Handle(ctx context.Context, command *CreateAccountCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createAccountHandler.Handle")
	defer span.Finish()

	passwordHashed, err := utils.HashPassword(command.Password)
	if err != nil {
		return err
	}

	account := &models.Account{
		AccountID:      command.AccountID,
		PlayerID:       command.PlayerID,
		Username:       command.Username,
		Email:          command.Email,
		PasswordHashed: passwordHashed,
		IsBan:          command.IsBan,
		CreatedAt:      command.CreatedAt,
		UpdatedAt:      command.UpdatedAt,
	}

	created, err := c.postgresRepo.CreateAccount(ctx, account)
	if err != nil {
		return err
	}

	c.redisRepo.PutAccount(ctx, created.AccountID.String(), created)
	return nil
}
