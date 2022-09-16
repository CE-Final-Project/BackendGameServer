package dto

type LoginAccountDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginAccountResponseDto struct {
	AccountID string `json:"account_id" validate:"required"`
	PlayerID  string `json:"player_id" validate:"required"`
	Token     string `json:"token" validate:"required"`
}
