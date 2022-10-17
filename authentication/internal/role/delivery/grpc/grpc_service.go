package grpc

import (
	"context"
	"database/sql"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/authentication/internal/role/commands"
	"github.com/ce-final-project/backend_game_server/authentication/internal/role/queries"
	"github.com/ce-final-project/backend_game_server/authentication/internal/role/service"
	GRPCService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/tracing"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
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

	var err error
	var id uint64
	md, _ := metadata.FromIncomingContext(ctx)

	id, err = strconv.ParseUint(md.Get("id")[0], 10, 64)
	s.log.Debug(id)
	if err != nil {
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	var role *models.Role
	command := commands.NewCreateRoleCommand(req.GetName(), id, id)
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	role, err = s.rs.Commands.CreateRole.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("CreateRole.Handle", err)
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, s.errResponse(codes.InvalidArgument, err)
		}
		return nil, s.errResponse(codes.Internal, err)
	}

	return &GRPCService.CreateRoleRes{Role: models.RoleToGrpcMessage(role)}, nil
}

func (s *grpcService) UpdateRole(ctx context.Context, req *GRPCService.UpdateRoleReq) (*GRPCService.UpdateRoleRes, error) {
	var span opentracing.Span
	ctx, span = tracing.StartGrpcServerTracerSpan(ctx, "RoleGRPCService.UpdateRole")
	defer span.Finish()

	md, _ := metadata.FromOutgoingContext(ctx)

	id, err := strconv.ParseUint(md.Copy().Get("id")[0], 10, 64)
	if err != nil {
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	var role *models.Role
	command := commands.NewChangeRoleNameCommand(req.ID, req.GetName(), id)
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	role, err = s.rs.Commands.ChangeRoleName.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("UpdateRole.Handle", err)
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, s.errResponse(codes.InvalidArgument, err)
		}
		return nil, s.errResponse(codes.Internal, err)
	}

	return &GRPCService.UpdateRoleRes{Role: models.RoleToGrpcMessage(role)}, nil
}

func (s *grpcService) SearchRole(ctx context.Context, req *GRPCService.SearchRolesReq) (*GRPCService.SearchRolesRes, error) {
	var span opentracing.Span
	ctx, span = tracing.StartGrpcServerTracerSpan(ctx, "RoleGRPCService.SearchRole")
	defer span.Finish()

	pg := utils.NewPaginationQuery(int(req.GetSize()), int(req.GetPage()))
	query := queries.NewSearchRoleQuery(req.GetSearch(), pg)
	if err := s.v.StructCtx(ctx, query); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	roles, err := s.rs.Queries.SearchRole.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetAccountByUsername.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return models.RoleListToGrpc(roles), nil
}

func (s *grpcService) DeleteRoleByID(ctx context.Context, req *GRPCService.DeleteRoleReq) (*GRPCService.DeleteRoleRes, error) {
	var span opentracing.Span
	ctx, span = tracing.StartGrpcServerTracerSpan(ctx, "RoleGRPCService.DeleteRoleByID")
	defer span.Finish()

	command := commands.NewDeleteRoleCommand(req.GetID())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	if err := s.rs.Commands.DeleteRole.Handle(ctx, command); err != nil {
		s.log.WarnMsg("DeleteRoleByID.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return &GRPCService.DeleteRoleRes{}, nil
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
