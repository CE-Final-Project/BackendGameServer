package queries

import (
	"context"
	accountService "github.com/ce-final-project/backend_game_server/account/proto"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type GetAccountByUsernameHandler interface {
	Handle(ctx context.Context, query *GetAccountByUsernameQuery) (*accountService.Account, error)
}

type getAccountByUsernameHandler struct {
	log      logger.Logger
	cfg      *config.Config
	asClient accountService.AccountServiceClient
}

func NewGetAccountByUsernameHandler(log logger.Logger, cfg *config.Config, asClient accountService.AccountServiceClient) *getAccountByUsernameHandler {
	return &getAccountByUsernameHandler{
		log:      log,
		cfg:      cfg,
		asClient: asClient,
	}
}

func (l *getAccountByUsernameHandler) Handle(ctx context.Context, query *GetAccountByUsernameQuery) (*accountService.Account, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetAccountByUsernameHandler.Handle")
	defer span.Finish()

	res, err := l.asClient.GetAccountByUsername(ctx, &accountService.GetAccountByUsernameReq{Username: query.Username})
	if err != nil {
		return nil, err
	}

	return res.GetAccount(), nil
}
