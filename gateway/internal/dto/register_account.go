package dto

type RegisterAccount struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email,max=320"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type RegisterAccountResponse struct {
	AccountID uint64 `json:"account_id" validate:"required"`
	PlayerID  string `json:"player_id" validate:"required"`
	Token     string `json:"token" validate:"required"`
}
