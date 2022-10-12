package commands

import (
	"context"
	"strconv"

	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"github.com/opentracing/opentracing-go"
)

type CreateAccountCmdHandler interface {
	Handle(ctx context.Context, command *CreateAccountCommand) error
}

type createAccountHandler struct {
	log         logger.Logger
	cfg         *config.Config
	accountRepo repository.AccountRepository
	cacheRepo   repository.CacheRepository
}

func NewCreateAccountHandler(log logger.Logger, cfg *config.Config, accountRepo repository.AccountRepository, cacheRepo repository.CacheRepository) CreateAccountCmdHandler {
	return &createAccountHandler{log: log, cfg: cfg, accountRepo: accountRepo, cacheRepo: cacheRepo}
}

func (c *createAccountHandler) Handle(ctx context.Context, command *CreateAccountCommand) error {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "createAccountHandler.Handle")
	defer span.Finish()

	passwordHashed, err := utils.HashPassword(command.Password)
	if err != nil {
		return err
	}

	account := &models.Account{
		Username:       command.Username,
		Email:          command.Email,
		PasswordHashed: passwordHashed,
	}
	account, err = c.accountRepo.InsertAccount(ctx, account)
	if err != nil {
		return err
	}

	c.cacheRepo.PutAccount(ctx, strconv.FormatUint(account.ID, 10), account)
	return nil
}
