package commands

type AccountCommands struct {
	CreateAccount  CreateAccountCmdHandler
	ChangePassword ChangePasswordCmdHandler
	DeleteAccount  DeleteAccountCmdHandler
}

func NewAccountCommands(
	createAccount CreateAccountCmdHandler,
	changePassword ChangePasswordCmdHandler,
	deleteAccount DeleteAccountCmdHandler,
) *AccountCommands {
	return &AccountCommands{CreateAccount: createAccount, ChangePassword: changePassword, DeleteAccount: deleteAccount}
}

type CreateAccountCommand struct {
	PlayerID string `json:"player_id,omitempty" validate:"required,max=15"`
	Username string `json:"username,omitempty" validate:"required,min=3,max=250"`
	Email    string `json:"email,omitempty" validate:"required,email,max=320"`
	Password string `json:"password,omitempty" validate:"required,min=8,max=100"`
}

func NewCreateAccountCommand(playerID, username, email, password string) *CreateAccountCommand {
	return &CreateAccountCommand{
		PlayerID: playerID,
		Username: username,
		Email:    email,
		Password: password,
	}
}

type ChangePasswordCommand struct {
	ID          uint64 `json:"id" validate:"required,numeric"`
	OldPassword string `json:"old_password,omitempty" bson:"old_password,omitempty" validate:"required"`
	NewPassword string `json:"new_password,omitempty" bson:"new_password,omitempty" validate:"required"`
}

func NewChangePasswordCommand(accountID uint64, oldPassword, newPassword string) *ChangePasswordCommand {
	return &ChangePasswordCommand{
		ID:          accountID,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}
}

type DeleteAccountCommand struct {
	ID uint64 `json:"id" validate:"required,numeric"`
}

func NewDeleteAccountCommand(accountID uint64) *DeleteAccountCommand {
	return &DeleteAccountCommand{ID: accountID}
}
