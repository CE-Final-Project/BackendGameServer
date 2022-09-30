package repository

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

type PostgresRepositoryMock struct {
	mock.Mock
}

func (p *PostgresRepositoryMock) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	args := p.Called(ctx, account)
	return args.Get(0).(*models.Account), args.Error(1)
}

func (p *PostgresRepositoryMock) UpdateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	args := p.Called(ctx, account)
	return args.Get(0).(*models.Account), args.Error(1)
}

func (p *PostgresRepositoryMock) DeleteAccount(ctx context.Context, uuid uuid.UUID) error {
	args := p.Called(ctx, uuid)
	return args.Error(0)
}

func (p *PostgresRepositoryMock) GetAccountById(ctx context.Context, uuid uuid.UUID) (*models.Account, error) {
	args := p.Called(ctx, uuid)
	return args.Get(0).(*models.Account), args.Error(1)
}

func (p *PostgresRepositoryMock) GetAccountByEmailOrUsername(ctx context.Context, email string, username string) (*models.Account, error) {
	args := p.Called(ctx, email, username)
	return args.Get(0).(*models.Account), args.Error(1)
}

func (p *PostgresRepositoryMock) Search(ctx context.Context, search string, pagination *utils.Pagination) (*models.AccountsList, error) {
	args := p.Called(ctx, search, pagination)
	return args.Get(0).(*models.AccountsList), args.Error(1)
}

func NewAccountRepositoryMock() *PostgresRepositoryMock {
	return &PostgresRepositoryMock{}
}
