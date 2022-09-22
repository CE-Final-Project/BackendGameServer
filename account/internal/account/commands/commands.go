package commands

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type AccountCommands struct {
	CreateAccount  CreateAccountCmdHandler
	UpdateAccount  UpdateAccountCmdHandler
	ChangePassword ChangePasswordCmdHandler
	BanAccountById BanAccountByIdCmdHandler
	DeleteAccount  DeleteAccountCmdHandler
}

func NewAccountCommands(
	createAccount CreateAccountCmdHandler,
	updateAccount UpdateAccountCmdHandler,
	changePassword ChangePasswordCmdHandler,
	banAccountById BanAccountByIdCmdHandler,
	deleteAccount DeleteAccountCmdHandler,
) *AccountCommands {
	return &AccountCommands{CreateAccount: createAccount, UpdateAccount: updateAccount, ChangePassword: changePassword, BanAccountById: banAccountById, DeleteAccount: deleteAccount}
}

type CreateAccountCommand struct {
	AccountID uuid.UUID `json:"account_id" bson:"_id,omitempty"`
	PlayerID  string    `json:"player_id,omitempty" bson:"player_id,omitempty" validate:"required,max=11"`
	Username  string    `json:"username,omitempty" bson:"username,omitempty" validate:"required,min=3,max=250"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty" validate:"required"`
	Password  string    `json:"password,omitempty" bson:"password,omitempty" validate:"required"`
	IsBan     bool      `json:"is_ban,omitempty" bson:"is_ban,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func NewCreateAccountCommand(accountID uuid.UUID, playerID, username, email, password string, isBan bool, createdAt, updatedAt time.Time) *CreateAccountCommand {
	return &CreateAccountCommand{
		AccountID: accountID,
		PlayerID:  playerID,
		Username:  username,
		Email:     email,
		Password:  password,
		IsBan:     isBan,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type UpdateAccountCommand struct {
	AccountID uuid.UUID `json:"account_id" bson:"_id,omitempty"`
	Username  string    `json:"username,omitempty" bson:"username,omitempty" validate:"required,min=3,max=250"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty" validate:"required"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func NewUpdateAccountCommand(accountID uuid.UUID, username, email string, updatedAt time.Time) *UpdateAccountCommand {
	return &UpdateAccountCommand{
		AccountID: accountID,
		Username:  username,
		Email:     email,
		UpdatedAt: updatedAt,
	}
}

type ChangePasswordCommand struct {
	AccountID   uuid.UUID `json:"account_id" bson:"account_id,omitempty"`
	OldPassword string    `json:"old_password,omitempty" bson:"old_password,omitempty" validate:"required"`
	NewPassword string    `json:"new_password,omitempty" bson:"new_password,omitempty" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func NewChangePasswordCommand(accountID uuid.UUID, oldPassword, newPassword string, updatedAt time.Time) *ChangePasswordCommand {
	return &ChangePasswordCommand{
		AccountID:   accountID,
		OldPassword: oldPassword,
		NewPassword: newPassword,
		UpdatedAt:   updatedAt,
	}
}

type BanAccountByIdCommand struct {
	AccountID uuid.UUID `json:"account_id" bson:"account_id,omitempty"`
	IsBan     bool      `json:"is_ban,omitempty" bson:"is_ban,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func NewBanAccountByIdCommand(accountID uuid.UUID, isBan bool, updatedAt time.Time) *BanAccountByIdCommand {
	return &BanAccountByIdCommand{AccountID: accountID, IsBan: isBan, UpdatedAt: updatedAt}
}

type DeleteAccountCommand struct {
	AccountID uuid.UUID `json:"account_id" bson:"_id,omitempty"`
}

func NewDeleteAccountCommand(accountID uuid.UUID) *DeleteAccountCommand {
	return &DeleteAccountCommand{AccountID: accountID}
}
