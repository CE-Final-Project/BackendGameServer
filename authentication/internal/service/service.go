package service

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/internal/dto"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/authentication/internal/repository"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/speps/go-hashids/v2"
	"math/rand"
	"time"
)

type AccountServices interface {
	CreateNewAccount(ctx context.Context, payload *dto.CreateAccountDto) (*dto.CreateAccountResponseDto, error)
	GetAccountById(ctx context.Context, accountID uuid.UUID) (*models.Account, error)
	GetAccountByUsernameOrEmail(ctx context.Context, username, email string) (*models.Account, error)
	UpdateAccount(ctx context.Context, payload *dto.UpdateAccountDto) (*dto.UpdateAccountResponseDto, error)
}

type AccountService struct {
	log       logger.Logger
	pgRepo    repository.Repository
	redisRepo repository.CacheRepository
}

func NewAccountService(log logger.Logger, pgRepo repository.Repository, redisRepo repository.CacheRepository) *AccountService {
	return &AccountService{
		log:       log,
		pgRepo:    pgRepo,
		redisRepo: redisRepo,
	}
}

func (a *AccountService) CreateNewAccount(ctx context.Context, payload *dto.CreateAccountDto) (*dto.CreateAccountResponseDto, error) {

	hd := hashids.NewData()
	hd.Salt = payload.Username + payload.Email
	hd.MinLength = 5
	h, err := hashids.NewWithData(hd)
	if err != nil {
		a.log.WarnMsg("hashIds.NewWithData", err)
		return nil, err
	}
	var playerID string
	playerID, err = h.Encode([]int{rand.Intn(5000)})
	if err != nil {
		a.log.WarnMsg("h.Encode", err)
		return nil, err
	}

	var hashedPassword string
	hashedPassword, err = utils.HashPassword(payload.Password)
	if err != nil {
		a.log.WarnMsg("utils.HashPassword", err)
		return nil, err
	}
	accountID := uuid.NewV4()
	account := &models.Account{
		AccountID:      accountID,
		PlayerID:       playerID,
		Username:       payload.Username,
		Email:          payload.Email,
		PasswordHashed: hashedPassword,
		IsBan:          false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	created, err := a.pgRepo.CreateAccount(ctx, account)
	a.log.Debugf("Created Account: %v", created)
	if err != nil {
		return nil, err
	}
	a.redisRepo.PutAccount(ctx, accountID.String(), account)

	return &dto.CreateAccountResponseDto{
		AccountID: created.AccountID,
		PlayerID:  created.PlayerID,
	}, nil
}

func (a *AccountService) GetAccountById(ctx context.Context, accountID uuid.UUID) (*models.Account, error) {
	if account, err := a.redisRepo.GetAccount(ctx, accountID.String()); err == nil && account != nil {
		return account, nil
	}
	account, err := a.pgRepo.GetAccountById(ctx, accountID)
	a.log.Debugf("Account Result: %v", account)
	if err != nil {
		return nil, err
	}
	a.redisRepo.PutAccount(ctx, account.AccountID.String(), account)
	return account, nil
}

func (a *AccountService) GetAccountByUsernameOrEmail(ctx context.Context, username, email string) (*models.Account, error) {

	account, err := a.pgRepo.GetAccountByEmailOrUsername(ctx, email, username)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (a *AccountService) UpdateAccount(ctx context.Context, payload *dto.UpdateAccountDto) (*dto.UpdateAccountResponseDto, error) {

	account := &models.Account{
		AccountID:      payload.AccountID,
		Username:       payload.Username,
		Email:          payload.Email,
		IsBan:          payload.IsBan,
		PasswordHashed: "",
	}
	if payload.Password != "" {
		hashedPassword, err := utils.HashPassword(payload.Password)
		if err != nil {
			a.log.WarnMsg("utils.HashPassword", err)
			return nil, err
		}
		account.PasswordHashed = hashedPassword
	}
	result, err := a.pgRepo.UpdateAccount(ctx, account)
	if err != nil {
		a.log.Errorf("can not update account: %v", err)
		return nil, err
	}
	a.redisRepo.DelAccount(ctx, result.AccountID.String())
	a.redisRepo.PutAccount(ctx, result.AccountID.String(), result)

	return &dto.UpdateAccountResponseDto{AccountID: result.AccountID, UpdatedAt: result.UpdatedAt}, nil
}

func (a *AccountService) ChangePasswordAccount(ctx context.Context, payload *dto.ChangePasswordDto) (*dto.UpdateAccountResponseDto, error) {
	hashedPassword, err := utils.HashPassword(payload.NewPassword)
	if err != nil {
		a.log.WarnMsg("utils.HashPassword", err)
		return nil, err
	}
	account := &models.Account{
		AccountID:      payload.AccountID,
		PasswordHashed: hashedPassword,
	}

	result, err := a.pgRepo.UpdateAccount(ctx, account)
	if err != nil {
		a.log.Errorf("can not update account: %v", err)
		return nil, err
	}

	a.redisRepo.DelAccount(ctx, result.AccountID.String())
	a.redisRepo.PutAccount(ctx, result.AccountID.String(), result)

	return &dto.UpdateAccountResponseDto{
		AccountID: result.AccountID,
		UpdatedAt: result.UpdatedAt,
	}, nil
}
