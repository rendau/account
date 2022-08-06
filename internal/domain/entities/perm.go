package entities

import (
	"github.com/rendau/dop/dopTypes"
)

type PermSt struct {
	Id       string `json:"id" db:"id"`
	App      string `json:"app" db:"app"`
	Dsc      string `json:"dsc" db:"dsc"`
	IsSystem bool   `json:"is_system" db:"is_system"`
}

type PermListSt struct {
	PermSt
}

type PermListParsSt struct {
	dopTypes.ListParams

	Ids    *[]string `json:"ids" form:"ids"`
	App    *string   `json:"app" form:"app"`
	Search *string   `json:"search" form:"search"`
}

type PermCUSt struct {
	Id  *string `json:"id" db:"id"`
	App *string `json:"app" db:"app"`
	Dsc *string `json:"dsc" db:"dsc"`
}
