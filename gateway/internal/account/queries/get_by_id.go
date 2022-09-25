package queries

import (
	"context"
	accountService "github.com/ce-final-project/backend_game_server/account/proto"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/dto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type GetAccountByIdHandler interface {
	Handle(ctx context.Context, query *GetAccountByIdQuery) (*dto.AccountResponseDto, error)
}

type getAccountByIdHandler struct {
	log      logger.Logger
	cfg      *config.Config
	asClient accountService.AccountServiceClient
}

func NewGetAccountByIdHandler(log logger.Logger, cfg *config.Config, asClient accountService.AccountServiceClient) *getAccountByIdHandler {
	return &getAccountByIdHandler{log: log, cfg: cfg, asClient: asClient}
}

func (q *getAccountByIdHandler) Handle(ctx context.Context, query *GetAccountByIdQuery) (*dto.AccountResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getAccountByIdHandler.Handle")
	defer span.Finish()

	res, err := q.asClient.GetAccountById(ctx, &accountService.GetAccountByIdReq{AccountID: query.AccountID.String()})
	if err != nil {
		return nil, err
	}

	return dto.AccountResponseFromGrpc(res.GetAccount()), nil
}
