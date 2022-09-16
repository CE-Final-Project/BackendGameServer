package service

import (
	"context"
	"github.com/ce-final-project/backend_game_server/gateway/internal/dto"
	uuid "github.com/satori/go.uuid"
)

type AuthService interface {
	Login(ctx context.Context, payload *dto.LoginAccountDto) (*dto.LoginAccountResponseDto, error)
	Register(ctx context.Context, payload *dto.RegisterAccountDto) (*dto.RegisterAccountResponseDto, error)
	VerifyToken(ctx context.Context, token string) (*dto.VerifyTokenResponseDto, error)
}

type AccountService interface {
	UpdateAccount(ctx context.Context, payload *dto.UpdateAccountDto) (*dto.UpdateAccountResponseDto, error)
	ChangePassword(ctx context.Context, payload *dto.ChangePasswordDto) (*dto.UpdateAccountResponseDto, error)
	GetAccountById(ctx context.Context, accountID uuid.UUID) (*dto.AccountResponseDto, error)
}

//type AuthService struct {
//	Commands *commands.AuthCommands
//	Queries  *queries.AuthQueries
//}
//
//func NewProductService(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer, asClient AuthService.AuthServiceClient) *AuthService {
//
//	registerAccountHandler := commands.NewRegisterAccountHandler(log, cfg, kafkaProducer)
//	updateAccountHandler := commands.NewUpdateAccountHandler(log, cfg, kafkaProducer)
//
//	getAccountByIdHandler := queries.NewGetAccountByIdHandler(log, cfg, asClient)
//
//	accountCommands := commands.NewAuthCommands(registerAccountHandler, updateAccountHandler)
//	accountQueries := queries.NewProductQueries(getAccountByIdHandler)
//
//	return &AuthService{Commands: accountCommands, Queries: accountQueries}
//}
