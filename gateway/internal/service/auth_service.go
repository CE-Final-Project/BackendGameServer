package service

import (
	"context"
	grpcAuthService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/dto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	uuid "github.com/satori/go.uuid"
)

type authService struct {
	log logger.Logger
	cfg *config.Config
	as  grpcAuthService.AuthServiceClient
}

func NewAuthService(log logger.Logger, cfg *config.Config, as grpcAuthService.AuthServiceClient) *authService {
	return &authService{
		log: log,
		cfg: cfg,
		as:  as,
	}
}

func (a *authService) Login(ctx context.Context, payload *dto.LoginAccountDto) (*dto.LoginAccountResponseDto, error) {
	result, err := a.as.Login(ctx, &grpcAuthService.LoginReq{
		Username: payload.Username,
		Password: payload.Password,
	})
	if err != nil {
		return nil, err
	}
	return &dto.LoginAccountResponseDto{
		AccountID: result.GetAccountID(),
		PlayerID:  result.GetPlayerID(),
		Token:     result.GetToken(),
	}, nil
}

func (a *authService) Register(ctx context.Context, payload *dto.RegisterAccountDto) (*dto.RegisterAccountResponseDto, error) {
	result, err := a.as.Register(ctx, &grpcAuthService.RegisterReq{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		return nil, err
	}
	AccountUUID, err := uuid.FromString(result.GetAccountID())
	if err != nil {
		return nil, err
	}

	return &dto.RegisterAccountResponseDto{
		AccountID: AccountUUID,
		PlayerID:  result.GetPlayerID(),
		Token:     result.GetToken(),
	}, nil
}

func (a *authService) VerifyToken(ctx context.Context, token string) (*dto.VerifyTokenResponseDto, error) {
	result, err := a.as.VerifyToken(ctx, &grpcAuthService.VerifyTokenReq{Token: token})
	if err != nil {
		return nil, err
	}
	AccountUUID, err := uuid.FromString(result.GetAccountID())
	if err != nil {
		return nil, err
	}
	return &dto.VerifyTokenResponseDto{
		Valid:     result.Valid,
		AccountID: AccountUUID,
		PlayerID:  result.GetPlayerID(),
	}, nil
}
