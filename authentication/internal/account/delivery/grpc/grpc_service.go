package grpc

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/account/commands"
	"github.com/ce-final-project/backend_game_server/authentication/internal/account/queries"
	"github.com/ce-final-project/backend_game_server/authentication/internal/account/service"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	accountService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/tracing"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcService struct {
	log logger.Logger
	cfg *config.Config
	v   *validator.Validate
	as  *service.AccountService
}

func (s *grpcService) GetAccountById(ctx context.Context, req *accountService.GetAccountByIdReq) (*accountService.GetAccountByIdRes, error) {
	var span opentracing.Span
	ctx, span = tracing.StartGrpcServerTracerSpan(ctx, "grpcService.GetAccountById")
	defer span.Finish()

	query := queries.NewGetAccountByIdQuery(req.GetAccountID())
	if err := s.v.StructCtx(ctx, query); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	account, err := s.as.Queries.GetAccountById.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetAccountById.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return &accountService.GetAccountByIdRes{Account: models.AccountToGrpcMessage(account)}, nil
}

func (s *grpcService) GetAccountByUsername(ctx context.Context, req *accountService.GetAccountByUsernameReq) (*accountService.GetAccountByUsernameRes, error) {
	var span opentracing.Span
	ctx, span = tracing.StartGrpcServerTracerSpan(ctx, "grpcService.GetAccountByUsername")
	defer span.Finish()

	query := queries.NewGetAccountByUsernameQuery(req.GetUsername())
	if err := s.v.StructCtx(ctx, query); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	account, err := s.as.Queries.GetAccountByUsername.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetAccountByUsername.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return &accountService.GetAccountByUsernameRes{Account: models.AccountToGrpcMessage(account)}, nil
}

func (s *grpcService) GetAccountByEmail(ctx context.Context, req *accountService.GetAccountByEmailReq) (*accountService.GetAccountByEmailRes, error) {
	var span opentracing.Span
	ctx, span = tracing.StartGrpcServerTracerSpan(ctx, "grpcService.GetAccountByEmail")
	defer span.Finish()

	query := queries.NewGetAccountByEmailQuery(req.GetEmail())
	if err := s.v.StructCtx(ctx, query); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	account, err := s.as.Queries.GetAccountByEmail.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetAccountByUsername.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return &accountService.GetAccountByEmailRes{Account: models.AccountToGrpcMessage(account)}, nil
}

func (s *grpcService) ChangePassword(ctx context.Context, req *accountService.ChangePasswordReq) (*accountService.ChangePasswordRes, error) {
	var span opentracing.Span
	ctx, span = tracing.StartGrpcServerTracerSpan(ctx, "grpcService.ChangePassword")
	defer span.Finish()

	command := commands.NewChangePasswordCommand(req.GetAccountID(), req.GetOldPassword(), req.GetNewPassword())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	err := s.as.Commands.ChangePassword.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("ChangPassword.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return &accountService.ChangePasswordRes{AccountID: command.ID}, nil
}

func (s *grpcService) SearchAccount(ctx context.Context, req *accountService.SearchAccountsReq) (*accountService.SearchAccountsRes, error) {
	var span opentracing.Span
	ctx, span = tracing.StartGrpcServerTracerSpan(ctx, "grpcService.SearchAccount")
	defer span.Finish()

	pq := utils.NewPaginationQuery(int(req.GetSize()), int(req.GetPage()))

	query := queries.NewSearchAccountQuery(req.GetSearch(), pq)
	if err := s.v.StructCtx(ctx, query); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	accounts, err := s.as.Queries.SearchAccount.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetAccountByUsername.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return models.AccountListToGrpc(accounts), nil
}

func (s *grpcService) DeleteAccountByID(ctx context.Context, req *accountService.DeleteAccountByIdReq) (*accountService.DeleteAccountByIdRes, error) {
	var span opentracing.Span
	ctx, span = tracing.StartGrpcServerTracerSpan(ctx, "grpcService.DeleteAccountByID")
	defer span.Finish()

	command := commands.NewDeleteAccountCommand(req.GetAccountID())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}
	err := s.as.Commands.DeleteAccount.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("DeleteAccountByID.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}
	return &accountService.DeleteAccountByIdRes{}, nil
}

func NewAccountGRPCService(log logger.Logger, cfg *config.Config, v *validator.Validate, as *service.AccountService) accountService.AccountServiceServer {
	return &grpcService{
		log: log,
		cfg: cfg,
		v:   v,
		as:  as,
	}
}

func (s *grpcService) errResponse(c codes.Code, err error) error {
	return status.Error(c, err.Error())
}
