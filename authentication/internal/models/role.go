package models

import (
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Role struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name" validate:"required,max=50"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	CreatedBy uint64    `json:"created_by,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UpdatedBy uint64    `json:"updated_by,omitempty"`
}

// RolesList roles list response with pagination
type RolesList struct {
	TotalCount int64   `json:"totalCount"`
	TotalPages int64   `json:"totalPages"`
	Page       int64   `json:"page"`
	Size       int64   `json:"size"`
	HasMore    bool    `json:"hasMore"`
	Roles      []*Role `json:"roles"`
}

func NewRoleListWithPagination(roles []*Role, count int64, pagination *utils.Pagination) *RolesList {
	return &RolesList{
		TotalCount: count,
		TotalPages: int64(pagination.GetTotalPages(int(count))),
		Page:       int64(pagination.GetPage()),
		Size:       int64(pagination.GetSize()),
		HasMore:    pagination.GetHasMore(int(count)),
		Roles:      roles,
	}
}

func RoleToGrpcMessage(role *Role) *authService.Role {
	return &authService.Role{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: timestamppb.New(role.CreatedAt),
		CreatedBy: role.CreatedBy,
		UpdatedAt: timestamppb.New(role.UpdatedAt),
		UpdatedBy: role.UpdatedBy,
	}
}

func RoleListToGrpc(roles *RolesList) *authService.SearchRolesRes {
	list := make([]*authService.Role, 0, len(roles.Roles))
	for _, account := range roles.Roles {
		list = append(list, RoleToGrpcMessage(account))
	}

	return &authService.SearchRolesRes{
		TotalCount: roles.TotalCount,
		TotalPages: roles.TotalPages,
		Page:       roles.Page,
		Size:       roles.Size,
		HasMore:    roles.HasMore,
		Roles:      list,
	}
}
