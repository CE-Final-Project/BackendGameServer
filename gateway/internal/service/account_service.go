package service

import (
	"context"
	grpcAuthService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/dto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	uuid "github.com/satori/go.uuid"
)

type accountService struct {
	log logger.Logger
	cfg *config.Config
	as  grpcAuthService.AuthServiceClient
}

func NewAccountService(log logger.Logger, cfg *config.Config, as grpcAuthService.AuthServiceClient) *accountService {
	return &accountService{
		log: log,
		cfg: cfg,
		as:  as,
	}
}

func (a *accountService) UpdateAccount(ctx context.Context, payload *dto.UpdateAccountDto) (*dto.UpdateAccountResponseDto, error) {

	accountReq := &grpcAuthService.UpdateAccountReq{
		AccountID: payload.AccountID.String(),
		Username:  payload.Username,
		Email:     payload.Email,
	}
	result, err := a.as.UpdateAccount(ctx, accountReq)
	if err != nil {
		return nil, err
	}
	return &dto.UpdateAccountResponseDto{
		AccountID: payload.AccountID,
		UpdatedAt: result.GetUpdatedAt().AsTime(),
	}, nil
}

func (a *accountService) ChangePassword(ctx context.Context, payload *dto.ChangePasswordDto) (*dto.UpdateAccountResponseDto, error) {
	changePassReq := &grpcAuthService.ChangePasswordReq{
		AccountID:   payload.AccountID.String(),
		OldPassword: payload.OldPassword,
		NewPassword: payload.NewPassword,
	}

	result, err := a.as.ChangePassword(ctx, changePassReq)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateAccountResponseDto{
		AccountID: payload.AccountID,
		UpdatedAt: result.GetUpdatedAt().AsTime(),
	}, nil

}

func (a *accountService) GetAccountById(ctx context.Context, accountID uuid.UUID) (*dto.AccountResponseDto, error) {
	result, err := a.as.GetAccountById(ctx, &grpcAuthService.GetAccountByIdReq{AccountID: accountID.String()})
	if err != nil {
		return nil, err
	}
	return &dto.AccountResponseDto{
		AccountID: result.GetAccountID(),
		PlayerID:  result.GetPlayerID(),
		Username:  result.GetUsername(),
		Email:     result.GetEmail(),
		IsBan:     result.GetIsBan(),
		CreatedAt: result.GetCreatedAt().AsTime(),
		UpdatedAt: result.GetUpdatedAt().AsTime(),
	}, nil
}
