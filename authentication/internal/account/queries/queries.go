package queries

import (
	"github.com/ce-final-project/backend_game_server/pkg/utils"
)

type AccountQueries struct {
	GetAccountById       GetAccountByIdHandler
	GetAccountByUsername GetAccountByUsernameHandler
	GetAccountByEmail    GetAccountByEmailHandler
	SearchAccount        SearchAccountHandler
}

func NewAccountQueries(getAccountById GetAccountByIdHandler, getAccountByUsername GetAccountByUsernameHandler, getAccountByEmail GetAccountByEmailHandler, searchAccount SearchAccountHandler) *AccountQueries {
	return &AccountQueries{
		GetAccountById:       getAccountById,
		GetAccountByUsername: getAccountByUsername,
		GetAccountByEmail:    getAccountByEmail,
		SearchAccount:        searchAccount,
	}
}

type GetAccountByIdQuery struct {
	ID uint64 `json:"id" validate:"required"`
}

func NewGetAccountByIdQuery(ID uint64) *GetAccountByIdQuery {
	return &GetAccountByIdQuery{ID: ID}
}

type GetAccountByUsernameQuery struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}

func NewGetAccountByUsernameQuery(username string) *GetAccountByUsernameQuery {
	return &GetAccountByUsernameQuery{Username: username}
}

type GetAccountByEmailQuery struct {
	Email string `json:"email" validate:"required,email,max=320"`
}

func NewGetAccountByEmailQuery(email string) *GetAccountByEmailQuery {
	return &GetAccountByEmailQuery{Email: email}
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
