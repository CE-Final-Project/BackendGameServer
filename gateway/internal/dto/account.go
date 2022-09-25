package dto

import (
	accountService "github.com/ce-final-project/backend_game_server/account/proto"
	"time"
)

type AccountResponseDto struct {
	AccountID string    `json:"account_id" validate:"required"`
	PlayerID  string    `json:"player_id" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	IsBan     bool      `json:"is_ban" validate:"required,boolean"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func AccountResponseFromGrpc(account *accountService.Account) *AccountResponseDto {
	return &AccountResponseDto{
		AccountID: account.GetAccountID(),
		PlayerID:  account.GetPlayerID(),
		Username:  account.GetUsername(),
		Email:     account.GetEmail(),
		IsBan:     account.GetIsBan(),
		CreatedAt: account.GetCreatedAt().AsTime(),
		UpdatedAt: account.GetUpdatedAt().AsTime(),
	}
}

type AccountsListResponseDto struct {
	TotalCount int64                 `json:"totalCount" bson:"totalCount"`
	TotalPages int64                 `json:"totalPages" bson:"totalPages"`
	Page       int64                 `json:"page" bson:"page"`
	Size       int64                 `json:"size" bson:"size"`
	HasMore    bool                  `json:"hasMore" bson:"hasMore"`
	Accounts   []*AccountResponseDto `json:"accounts" bson:"accounts"`
}

func AccountsListResponseFromGrpc(listResponse *accountService.SearchRes) *AccountsListResponseDto {
	list := make([]*AccountResponseDto, 0, len(listResponse.GetAccounts()))
	for _, product := range listResponse.GetAccounts() {
		list = append(list, AccountResponseFromGrpc(product))
	}

	return &AccountsListResponseDto{
		TotalCount: listResponse.GetTotalCount(),
		TotalPages: listResponse.GetTotalPages(),
		Page:       listResponse.GetPage(),
		Size:       listResponse.GetSize(),
		HasMore:    listResponse.GetHasMore(),
		Accounts:   list,
	}
}
