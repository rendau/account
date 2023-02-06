package entities

import (
	"github.com/rendau/dop/dopTypes"
)

type RoleSt struct {
	Id       int64  `json:"id" db:"id"`
	Code     string `json:"code" db:"code"`
	Name     string `json:"name" db:"name"`
	IsSystem bool   `json:"is_system" db:"is_system"`

	PermIds []int64 `json:"perm_ids" db:"perm_ids"`
}

type RoleListSt struct {
	Id       int64  `json:"id" db:"id"`
	Code     string `json:"code" db:"code"`
	Name     string `json:"name" db:"name"`
	IsSystem bool   `json:"is_system" db:"is_system"`
}

type RoleListParsSt struct {
	dopTypes.ListParams

	Ids *[]int64 `json:"ids" form:"ids"`
}

type RoleCUSt struct {
	Code *string `json:"code" db:"code"`
	Name *string `json:"name" db:"name"`

	PermIds []int64 `json:"perm_ids"`
}
