package queries

import (
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	uuid "github.com/satori/go.uuid"
)

type AccountQueries struct {
	GetAccountById       GetAccountByIdHandler
	GetAccountByUsername GetAccountByUsernameHandler
	SearchAccount        SearchAccountHandler
}

func NewAccountQueries(getAccountById GetAccountByIdHandler, getAccountByUsername GetAccountByUsernameHandler, searchAccount SearchAccountHandler) *AccountQueries {
	return &AccountQueries{GetAccountById: getAccountById, GetAccountByUsername: getAccountByUsername, SearchAccount: searchAccount}
}

type GetAccountByIdQuery struct {
	AccountID uuid.UUID `json:"account_id" validate:"required,gte=0,lte=255"`
}

func NewGetAccountByIdQuery(accountID uuid.UUID) *GetAccountByIdQuery {
	return &GetAccountByIdQuery{AccountID: accountID}
}

type GetAccountByUsernameQuery struct {
	Username string `json:"username" validate:"required"`
}

func NewGetAccountByUsernameQuery(username string) *GetAccountByUsernameQuery {
	return &GetAccountByUsernameQuery{Username: username}
}

type SearchAccountQuery struct {
	Text       string            `json:"text"`
	Pagination *utils.Pagination `json:"pagination"`
}

func NewSearchAccountQuery(text string, pagination *utils.Pagination) *SearchAccountQuery {
	return &SearchAccountQuery{
		Text:       text,
		Pagination: pagination,
	}
}