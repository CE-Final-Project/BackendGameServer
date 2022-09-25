package dto

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type ChangePassword struct {
	AccountID   uuid.UUID `json:"account_id" validate:"required"`
	OldPassword string    `json:"old_password" validate:"required"`
	NewPassword string    `json:"new_password" validate:"required"`
}

type ChangePasswordResponse struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
}
