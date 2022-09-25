package server

import (
	"context"
	accountService "github.com/ce-final-project/backend_game_server/account/proto"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/auth/service"
	"github.com/ce-final-project/backend_game_server/authentication/internal/client"
	kafkaClient "github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	log logger.Logger
	cfg *config.Config
	v   *validator.Validate
	as  *service.AuthService
}

func NewServer(log logger.Logger, cfg *config.Config) *server {
	return &server{
		log: log,
		v:   validator.New(),
		cfg: cfg,
	}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	accountServiceConn, err := client.NewAccountServiceConn(ctx, s.cfg)
	if err != nil {
		return err
	}
	defer func(accountServiceConn *grpc.ClientConn) {
		err := accountServiceConn.Close()
		if err != nil {
			s.log.WarnMsg("close accountServiceConn", err)
		}
	}(accountServiceConn)
	asClient := accountService.NewAccountServiceClient(accountServiceConn)

	kafkaProducer := kafkaClient.NewProducer(s.log, s.cfg.Kafka.Brokers)
	defer kafkaProducer.Close()

	s.as = service.NewAuthService(s.log, s.cfg, kafkaProducer, asClient)

	closeGrpcServer, grpcServer, err := s.NewAuthGrpcServer()
	if err != nil {
		return errors.Wrap(err, "NewScmGrpcServer")
	}
	defer closeGrpcServer()

	<-ctx.Done()
	grpcServer.GracefulStop()
	return nil
}
