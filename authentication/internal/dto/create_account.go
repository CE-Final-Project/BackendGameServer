package dto

import uuid "github.com/satori/go.uuid"

type CreateAccountDto struct {
	Username string `json:"username,omitempty" bson:"username,omitempty" validate:"required,min=3,max=250"`
	Email    string `json:"email,omitempty" bson:"email,omitempty" validate:"required"`
	Password string `json:"password_hash,omitempty" bson:"password_hash,omitempty" validate:"required"`
}

type CreateAccountResponseDto struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	PlayerID  string    `json:"player_id" validate:"required,max=11"`
}
