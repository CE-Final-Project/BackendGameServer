package queries

type AuthQueries struct {
	GetAccountByUsername GetAccountByUsernameHandler
}

func NewAuthQueries(getAccountByUsername GetAccountByUsernameHandler) *AuthQueries {
	return &AuthQueries{
		GetAccountByUsername: getAccountByUsername,
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
