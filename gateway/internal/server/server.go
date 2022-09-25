package server

import (
	"context"
	accountService "github.com/ce-final-project/backend_game_server/account/proto"
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	accountClient "github.com/ce-final-project/backend_game_server/gateway/internal/account/client"
	accountV1 "github.com/ce-final-project/backend_game_server/gateway/internal/account/delivery/http/v1"
	sAcc "github.com/ce-final-project/backend_game_server/gateway/internal/account/service"
	authClient "github.com/ce-final-project/backend_game_server/gateway/internal/auth/client"
	authV1 "github.com/ce-final-project/backend_game_server/gateway/internal/auth/delivery/http/v1"
	sAuth "github.com/ce-final-project/backend_game_server/gateway/internal/auth/service"
	"github.com/ce-final-project/backend_game_server/gateway/internal/middlewares"
	"github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	log  logger.Logger
	cfg  *config.Config
	v    *validator.Validate
	mw   middlewares.MiddlewareManager
	echo *echo.Echo
	acs  *sAcc.AccountService
	as   *sAuth.AuthService
}

func NewServer(log logger.Logger, cfg *config.Config) *server {
	return &server{
		log:  log,
		cfg:  cfg,
		v:    validator.New(),
		echo: echo.New(),
	}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	authServiceConn, err := authClient.NewAuthServiceConn(ctx, s.cfg)
	if err != nil {
		return err
	}
	defer authServiceConn.Close()
	asClient := authService.NewAuthServiceClient(authServiceConn)

	accountServiceConn, err := accountClient.NewAccountServiceConn(ctx, s.cfg)
	if err != nil {
		return err
	}
	defer accountServiceConn.Close()

	acsClient := accountService.NewAccountServiceClient(accountServiceConn)

	kafkaProducer := kafka.NewProducer(s.log, s.cfg.Kafka.Brokers)

	s.acs = sAcc.NewAccountService(s.log, s.cfg, kafkaProducer, acsClient)
	s.as = sAuth.NewAuthService(s.log, s.cfg, kafkaProducer, asClient)
	s.mw = middlewares.NewMiddlewareManager(s.log, s.cfg, s.as)

	authHandler := authV1.NewAuthsHandlers(s.echo.Group("/api/v1"), s.log, s.mw, s.cfg, s.as, s.v)
	authHandler.MapRoutes()

	accountHandler := accountV1.NewAccountHandlers(s.echo.Group("/api/v1"), s.log, s.mw, s.cfg, s.acs, s.v)
	accountHandler.MapRoutes()

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
