package queries

import (
	"context"
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/dto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type VerifyTokenHandler interface {
	Handle(ctx context.Context, query *VerifyTokenQuery) (*dto.VerifyTokenResponse, error)
}

type verifyTokenHandler struct {
	log      logger.Logger
	cfg      *config.Config
	asClient authService.AuthServiceClient
}

func NewVerifyTokenHandler(log logger.Logger, cfg *config.Config, asClient authService.AuthServiceClient) *verifyTokenHandler {
	return &verifyTokenHandler{log: log, cfg: cfg, asClient: asClient}
}

func (q *verifyTokenHandler) Handle(ctx context.Context, query *VerifyTokenQuery) (*dto.VerifyTokenResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "verifyTokenHandler.Handle")
	defer span.Finish()

	res, err := q.asClient.VerifyToken(ctx, &authService.VerifyTokenReq{Token: query.Token})
	if err != nil {
		return nil, err
	}

	return dto.VerifyTokenResponseFromGRPC(res)
}
