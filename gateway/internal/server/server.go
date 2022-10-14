package server

import (
	"context"
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	authClient "github.com/ce-final-project/backend_game_server/gateway/internal/auth/client"
	authV1 "github.com/ce-final-project/backend_game_server/gateway/internal/auth/delivery/http/v1"
	sAuth "github.com/ce-final-project/backend_game_server/gateway/internal/auth/service"
	"github.com/ce-final-project/backend_game_server/gateway/internal/middlewares"
	"github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	log  logger.Logger
	cfg  *config.Config
	v    *validator.Validate
	mw   middlewares.MiddlewareManager
	echo *echo.Echo
	as   *sAuth.AuthService
}

func NewServer(log logger.Logger, cfg *config.Config) *Server {
	return &Server{
		log:  log,
		cfg:  cfg,
		v:    validator.New(),
		echo: echo.New(),
	}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	authServiceConn, err := authClient.NewAuthServiceConn(ctx, s.cfg)
	if err != nil {
		return err
	}
	defer func(authServiceConn *grpc.ClientConn) {
		err := authServiceConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(authServiceConn)
	asClient := authService.NewAuthServiceClient(authServiceConn)
	accClient := authService.NewAccountServiceClient(authServiceConn)

	kafkaProducer := kafka.NewProducer(s.log, s.cfg.Kafka.Brokers)

	s.as = sAuth.NewAuthService(s.log, s.cfg, kafkaProducer, asClient, accClient)
	s.mw = middlewares.NewMiddlewareManager(s.log, s.cfg, s.as)

	authHandler := authV1.NewAuthsHandlers(s.echo.Group("/api/v1"), s.log, s.mw, s.cfg, s.as, s.v)
	authHandler.MapRoutes()

	go func() {
		if err := s.runHttpServer(); err != nil {
			s.log.Errorf(" s.runHttpServer: %v", err)
			cancel()
		}
	}()
	s.log.Infof("API Gateway is listening on PORT: %s", s.cfg.HTTP.Port)
	<-ctx.Done()
	if err := s.echo.Server.Shutdown(ctx); err != nil {
		s.log.WarnMsg("echo.Server.Shutdown", err)
	}
	return nil
}
