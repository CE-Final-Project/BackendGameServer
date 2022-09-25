package queries

type AuthQueries struct {
	VerifyToken VerifyTokenHandler
}

func NewAccountQueries(verifyToken VerifyTokenHandler) *AuthQueries {
	return &AuthQueries{VerifyToken: verifyToken}
}

type VerifyTokenQuery struct {
	Token string `json:"token" validate:"required"`
}

func NewVerifyTokenQuery(token string) *VerifyTokenQuery {
	return &VerifyTokenQuery{Token: token}
}
