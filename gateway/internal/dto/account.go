package dto

type Account struct {
	ID       uint64 `json:"id" validate:"required,numeric"`
	PlayerID string `json:"player_id" validate:"required,max=15"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
}
