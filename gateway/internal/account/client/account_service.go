package client

import (
	"context"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const (
	backoffLinear  = 100 * time.Millisecond
	backoffRetries = 3
)

func NewAccountServiceConn(ctx context.Context, cfg *config.Config) (*grpc.ClientConn, error) {

	opts := []grpcRetry.CallOption{
		grpcRetry.WithBackoff(grpcRetry.BackoffLinear(backoffLinear)),
		grpcRetry.WithCodes(codes.NotFound, codes.Aborted),
		grpcRetry.WithMax(backoffRetries),
	}

	accountServiceConn, err := grpc.DialContext(
		ctx,
		cfg.GRPC.AccountServicePort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpcRetry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		return nil, errors.Wrap(err, "grpc.DialContext")
	}

	return accountServiceConn, nil
}
