package repository

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/internal/models"
	"github.com/stretchr/testify/mock"
)

type AccountCacheRepositoryMock struct {
	mock.Mock
}

func (a *AccountCacheRepositoryMock) PutAccount(ctx context.Context, key string, account *models.Account) {
	_ = a.Called(ctx, key, account)
	return
}

func (a *AccountCacheRepositoryMock) PutKeyReference(ctx context.Context, key string, targetKey string) {
	_ = a.Called(ctx, key, targetKey)
	return
}

func (a *AccountCacheRepositoryMock) GetAccount(ctx context.Context, key string) (*models.Account, error) {
	args := a.Called(ctx, key)
	return args.Get(0).(*models.Account), args.Error(1)
}

func (a *AccountCacheRepositoryMock) GetAccountReference(ctx context.Context, key string) (*models.Account, error) {
	args := a.Called(ctx, key)
	return args.Get(0).(*models.Account), args.Error(1)
}

func (a *AccountCacheRepositoryMock) DelAccount(ctx context.Context, key string) {
	_ = a.Called(ctx, key)
}

func (a *AccountCacheRepositoryMock) DelAllAccounts(ctx context.Context) {
	_ = a.Called(ctx)
}

func NewAccountCacheRepositoryMock() *AccountCacheRepositoryMock {
	return &AccountCacheRepositoryMock{}
}
