package dto

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type UpdateAccountDto struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	Username  string    `json:"username,omitempty" bson:"username,omitempty" validate:"required,min=3,max=250"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty" validate:"required"`
	Password  string    `json:"password,omitempty" bson:"password,omitempty" validate:"required"`
	IsBan     bool      `json:"is_ban,omitempty" bson:"is_ban,omitempty"`
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
