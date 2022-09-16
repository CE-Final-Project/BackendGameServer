package client

import (
	"context"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const (
	backoffLinear  = 100 * time.Millisecond
	backoffRetries = 3
)

func NewAuthServiceConn(ctx context.Context, cfg *config.Config) (*grpc.ClientConn, error) {

	authServiceConn, err := grpc.DialContext(
		ctx,
		cfg.GRPC.AuthServicePort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "grpc.DialContext")
	}

	return authServiceConn, nil
}
