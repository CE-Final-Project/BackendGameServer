package grpc

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/role/commands"
	"github.com/ce-final-project/backend_game_server/authentication/internal/role/service"
	GRPCService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/tracing"
	"github.com/go-playground/validator"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
)

type grpcService struct {
	log logger.Logger
	cfg *config.Config
	v   *validator.Validate
	rs  *service.RoleService
}

func (s *grpcService) CreateRole(ctx context.Context, req *GRPCService.CreateRoleReq) (*GRPCService.CreateRoleRes, error) {
	var span opentracing.Span
	ctx, span = tracing.StartGrpcServerTracerSpan(ctx, "RoleGRPCService.CreateRole")
	defer span.Finish()

	md, _ := metadata.FromOutgoingContext(ctx)

	id, err := strconv.ParseUint(md.Get("id")[0], 10, 64)
	if err != nil {
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	command := commands.NewCreateRoleCommand(req.GetName(), id, id)
	role, err :=
}

func (s *grpcService) UpdateRole(ctx context.Context, req *GRPCService.UpdateRoleReq) (*GRPCService.UpdateRoleRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s *grpcService) SearchRole(ctx context.Context, req *GRPCService.SearchRolesReq) (*GRPCService.SearchRolesRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s *grpcService) DeleteRoleByID(ctx context.Context, req *GRPCService.DeleteRoleReq) (*GRPCService.DeleteRoleRes, error) {
	//TODO implement me
	panic("implement me")
}

func NewRoleGRPCService(log logger.Logger, cfg *config.Config, v *validator.Validate, rs *service.RoleService) GRPCService.RoleServiceServer {
	return &grpcService{
		log: log,
		cfg: cfg,
		v:   v,
		rs:  rs,
	}
}

func (s *grpcService) errResponse(c codes.Code, err error) error {
	return status.Error(c, err.Error())
}
