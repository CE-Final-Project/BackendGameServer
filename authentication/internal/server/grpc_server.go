package server

import (
	grpcAccountDelivery "github.com/ce-final-project/backend_game_server/authentication/internal/account/delivery/grpc"
	grpcAuthDelivery "github.com/ce-final-project/backend_game_server/authentication/internal/auth/delivery/grpc"
	grpcRoleDelivery "github.com/ce-final-project/backend_game_server/authentication/internal/role/delivery/grpc"
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

func (s *Server) NewAuthGrpcServer() (func() error, *grpc.Server, error) {
	l, err := net.Listen("tcp", s.cfg.GRPC.Port)
	if err != nil {
		return nil, nil, errors.Wrap(err, "net.Listen")
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
	)

	grpcAuthService := grpcAuthDelivery.NewAuthGRPCService(s.log, s.cfg, s.v, s.acc)
	authService.RegisterAuthServiceServer(grpcServer, grpcAuthService)

	grpcAccountService := grpcAccountDelivery.NewAccountGRPCService(s.log, s.cfg, s.v, s.acc)
	authService.RegisterAccountServiceServer(grpcServer, grpcAccountService)

	grpcRoleService := grpcRoleDelivery.NewRoleGRPCService(s.log, s.cfg, s.v, s.rs)
	authService.RegisterRoleServiceServer(grpcServer, grpcRoleService)

	if s.cfg.GRPC.Development {
		reflection.Register(grpcServer)
	}

	go func() {
		s.log.Infof("auth gRPC server is listening on port: %s", s.cfg.GRPC.Port)
		s.log.Fatal(grpcServer.Serve(l))
	}()

	return l.Close, grpcServer, nil
}
