package dto

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type UpdateAccountDto struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty" validate:"email"`
}

type UpdateAccountResponseDto struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChangePasswordDto struct {
	AccountID   uuid.UUID `json:"account_id" validate:"required"`
	OldPassword string    `json:"old_password" validate:"required"`
	NewPassword string    `json:"new_password" validate:"required"`
}
