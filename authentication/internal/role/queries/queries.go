package queries

import "github.com/ce-final-project/backend_game_server/pkg/utils"

type RoleQueries struct {
	SearchRole    SearchRoleQueryHandler
	GetRoleByID   GetRoleByIDQueryHandler
	GetRoleByName GetRoleByNameQueryHandler
}

func NewRoleQueries(searchRole SearchRoleQueryHandler, getRoleByID GetRoleByIDQueryHandler, getRoleByName GetRoleByNameQueryHandler) *RoleQueries {
	return &RoleQueries{
		SearchRole:    searchRole,
		GetRoleByID:   getRoleByID,
		GetRoleByName: getRoleByName,
	}
}

type SearchRoleQuery struct {
	Text       string            `json:"text"`
	Pagination *utils.Pagination `json:"pagination"`
}

func NewSearchRoleQuery(text string, pagination *utils.Pagination) *SearchRoleQuery {
	return &SearchRoleQuery{
		Text:       text,
		Pagination: pagination,
	}
}

type GetRoleByIDQuery struct {
	ID uint64 `json:"id" validate:"required,numeric"`
}

func NewGetRoleByIDQuery(roleID uint64) *GetRoleByIDQuery {
	return &GetRoleByIDQuery{ID: roleID}
}

type GetRoleByNameQuery struct {
	Name string `json:"name" validate:"required,max=50"`
}

func NewGetRoleByNameQuery(name string) *GetRoleByNameQuery {
	return &GetRoleByNameQuery{Name: name}
}
