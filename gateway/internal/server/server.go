package server

import (
	"context"
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/client"
	v1 "github.com/ce-final-project/backend_game_server/gateway/internal/delivery/http/v1"
	"github.com/ce-final-project/backend_game_server/gateway/internal/middlewares"
	"github.com/ce-final-project/backend_game_server/gateway/internal/service"
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
	acs  service.AccountService
	as   service.AuthService
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

	authServiceConn, err := client.NewAuthServiceConn(ctx, s.cfg)
	if err != nil {
		return err
	}
	defer authServiceConn.Close()
	asClient := authService.NewAuthServiceClient(authServiceConn)

	s.acs = service.NewAccountService(s.log, s.cfg, asClient)
	s.as = service.NewAuthService(s.log, s.cfg, asClient)
	s.mw = middlewares.NewMiddlewareManager(s.log, s.cfg, s.as)

	authHandler := v1.NewAuthsHandlers(s.echo.Group("/api/v1"), s.log, s.mw, s.cfg, s.acs, s.as, s.v)
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
