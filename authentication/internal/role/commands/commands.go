package commands

type RoleCommands struct {
	createRole     CreateRoleCmdHandler
	changeRoleName ChangeRoleNameCmdHandler
	deleteRole     DeleteRoleCmdHandler
}

func NewRoleCommands(createRole CreateRoleCmdHandler, changeRoleName ChangeRoleNameCmdHandler, deleteRole DeleteRoleCmdHandler) *RoleCommands {
	return &RoleCommands{
		createRole:     createRole,
		changeRoleName: changeRoleName,
		deleteRole:     deleteRole,
	}
}

type CreateRoleCommand struct {
	Name      string `json:"name" validate:"required,max=50"`
	CreatedBy uint64 `json:"created_by" validate:"required,numeric"`
	UpdatedBy uint64 `json:"updated_by" validate:"required,numeric"`
}

func NewCreateRoleCommand(name string, createdBy, updatedBy uint64) *CreateRoleCommand {
	return &CreateRoleCommand{
		Name:      name,
		CreatedBy: createdBy,
		UpdatedBy: updatedBy,
	}
}

type ChangeRoleNameCommand struct {
	ID        uint64 `json:"id" validate:"required,numeric"`
	Name      string `json:"name" validate:"required,max=50"`
	UpdatedBy uint64 `json:"updated_by" validate:"required,numeric"`
}

func NewChangeRoleNameCommand(roleID uint64, name string, updatedBy uint64) *ChangeRoleNameCommand {
	return &ChangeRoleNameCommand{
		ID:        roleID,
		Name:      name,
		UpdatedBy: updatedBy,
	}
}

type DeleteRoleCommand struct {
	ID uint64 `json:"id" validate:"required,numeric"`
}

func NewDeleteRoleCommand(roleID uint64) *DeleteRoleCommand {
	return &DeleteRoleCommand{ID: roleID}
}
