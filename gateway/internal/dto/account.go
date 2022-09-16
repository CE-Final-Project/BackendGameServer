package dto

import "time"

type AccountResponseDto struct {
	AccountID string    `json:"account_id" validate:"required"`
	PlayerID  string    `json:"player_id" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	IsBan     bool      `json:"is_ban" validate:"required,boolean"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
