package models

import (
	"github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Account struct {
	ID             uint64    `json:"id"`
	PlayerID       string    `json:"player_id,omitempty" validate:"required,max=11"`
	Username       string    `json:"username,omitempty" validate:"required,min=3,max=250"`
	Email          string    `json:"email,omitempty" validate:"required,email,max=320"`
	PasswordHashed string    `json:"password_hashed,omitempty" validate:"required,max=255"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

// AccountsList accounts list response with pagination
type AccountsList struct {
	TotalCount int64      `json:"totalCount" bson:"totalCount"`
	TotalPages int64      `json:"totalPages" bson:"totalPages"`
	Page       int64      `json:"page" bson:"page"`
	Size       int64      `json:"size" bson:"size"`
	HasMore    bool       `json:"hasMore" bson:"hasMore"`
	Accounts   []*Account `json:"accounts" bson:"accounts"`
}

func NewAccountListWithPagination(accounts []*Account, count int64, pagination *utils.Pagination) *AccountsList {
	return &AccountsList{
		TotalCount: count,
		TotalPages: int64(pagination.GetTotalPages(int(count))),
		Page:       int64(pagination.GetPage()),
		Size:       int64(pagination.GetSize()),
		HasMore:    pagination.GetHasMore(int(count)),
		Accounts:   accounts,
	}
}

func AccountToGrpcMessage(account *Account) *authService.Account {
	return &authService.Account{
		AccountID: account.ID,
		PlayerID:  account.PlayerID,
		Username:  account.Username,
		Email:     account.Email,
		CreatedAt: timestamppb.New(account.CreatedAt),
		UpdatedAt: timestamppb.New(account.UpdatedAt),
	}
}

func AccountListToGrpc(accounts *AccountsList) *authService.SearchAccountsRes {
	list := make([]*authService.Account, 0, len(accounts.Accounts))
	for _, account := range accounts.Accounts {
		list = append(list, AccountToGrpcMessage(account))
	}

	return &authService.SearchAccountsRes{
		TotalCount: accounts.TotalCount,
		TotalPages: accounts.TotalPages,
		Page:       accounts.Page,
		Size:       accounts.Size,
		HasMore:    accounts.HasMore,
		Accounts:   list,
	}
}
