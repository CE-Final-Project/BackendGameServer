package grpc

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/config"
	"github.com/ce-final-project/backend_game_server/account/internal/account/commands"
	"github.com/ce-final-project/backend_game_server/account/internal/account/queries"
	"github.com/ce-final-project/backend_game_server/account/internal/account/service"
	"github.com/ce-final-project/backend_game_server/account/internal/models"
	accountService "github.com/ce-final-project/backend_game_server/account/proto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"github.com/ce-final-project/backend_rest_api/pkg/tracing"
	"github.com/go-playground/validator"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type grpcService struct {
	log logger.Logger
	cfg *config.Config
	v   *validator.Validate
	as  *service.AccountService
}

func NewGrpcService(log logger.Logger, cfg *config.Config, v *validator.Validate, as *service.AccountService) *grpcService {
	return &grpcService{
		log: log,
		cfg: cfg,
		v:   v,
		as:  as,
	}
}

func (s *grpcService) CreateAccount(ctx context.Context, req *accountService.CreateAccountReq) (*accountService.CreateAccountRes, error) {

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.CreateAccount")
	defer span.Finish()

	command := commands.NewCreateAccountCommand(req.GetAccountID(), req.GetPlayerID(), req.GetUsername(), req.GetEmail(), req.GetPassword(), false, time.Now(), time.Now())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	if err := s.as.Commands.CreateAccount.Handle(ctx, command); err != nil {
		s.log.WarnMsg("CreateAccount.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	return &accountService.CreateAccountRes{AccountID: req.GetAccountID()}, nil
}

func (s *grpcService) UpdateAccount(ctx context.Context, req *accountService.UpdateAccountReq) (*accountService.UpdateAccountRes, error) {

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.UpdateAccount")
	defer span.Finish()

	command := commands.NewUpdateAccountCommand(req.GetAccountID(), req.GetUsername(), req.GetEmail(), time.Now())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	if err := s.as.Commands.UpdateAccount.Handle(ctx, command); err != nil {
		s.log.WarnMsg("UpdateAccount.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	return &accountService.UpdateAccountRes{AccountID: req.GetAccountID()}, nil
}

func (s *grpcService) GetAccountById(ctx context.Context, req *accountService.GetAccountByIdReq) (*accountService.GetAccountByIdRes, error) {

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.GetAccountById")
	defer span.Finish()

	accountUUID, err := uuid.FromString(req.GetAccountID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	query := queries.NewGetAccountByIdQuery(accountUUID)
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

func (s *grpcService) SearchAccount(ctx context.Context, req *accountService.SearchReq) (*accountService.SearchRes, error) {

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.SearchAccount")
	defer span.Finish()

	pq := utils.NewPaginationQuery(int(req.GetSize()), int(req.GetPage()))

	query := queries.NewSearchAccountQuery(req.GetSearch(), pq)
	accountsList, err := s.as.Queries.SearchAccount.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("SearchAccount.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return models.AccountListToGrpc(accountsList), nil
}

func (s *grpcService) DeleteAccountByID(ctx context.Context, req *accountService.DeleteAccountByIdReq) (*accountService.DeleteAccountByIdRes, error) {

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.DeleteAccountByID")
	defer span.Finish()

	accountUUID, err := uuid.FromString(req.GetAccountID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	if err := s.as.Commands.DeleteAccount.Handle(ctx, commands.NewDeleteAccountCommand(accountUUID)); err != nil {
		s.log.WarnMsg("DeleteAccount.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return &accountService.DeleteAccountByIdRes{}, nil
}

func (s *grpcService) errResponse(c codes.Code, err error) error {
	return status.Error(c, err.Error())
}
