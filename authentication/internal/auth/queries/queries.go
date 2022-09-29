package queries

type AuthQueries struct {
	GetAccountByUsername GetAccountByUsernameHandler
	GetAccountById       GetAccountByIdHandler
}

func NewAuthQueries(getAccountByUsername GetAccountByUsernameHandler, getAccountById GetAccountByIdHandler) *AuthQueries {
	return &AuthQueries{
		GetAccountByUsername: getAccountByUsername,
		GetAccountById:       getAccountById,
	}
}

type GetAccountByUsernameQuery struct {
	Username string `json:"username,omitempty" validate:"required"`
}

func NewGetAccountByUsernameQuery(username string) *GetAccountByUsernameQuery {
	return &GetAccountByUsernameQuery{
		Username: username,
	}
}

type GetAccountByIdQuery struct {
	AccountID string `json:"account_id" validate:"required"`
}

func NewGetAccountByIdQuery(accountID string) *GetAccountByIdQuery {
	return &GetAccountByIdQuery{AccountID: accountID}
}
