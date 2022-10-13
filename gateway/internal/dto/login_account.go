package dto

type LoginAccount struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type LoginAccountResponse struct {
	AccountID uint64 `json:"account_id" validate:"required"`
	PlayerID  string `json:"player_id" validate:"required"`
	Token     string `json:"token" validate:"required"`
}
