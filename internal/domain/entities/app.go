package entities

import (
	"github.com/rendau/dop/dopTypes"
)

type AppSt struct {
	Id           int64  `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	PermPrefix   string `json:"perm_prefix" db:"perm_prefix"`
	PermUrl      string `json:"perm_url" db:"perm_url"`
	IsAccountApp bool   `json:"is_account_app" db:"is_account_app"`
}

type AppListParsSt struct {
	dopTypes.ListParams

	Ids  *[]int64 `json:"ids" form:"ids"`
	Name *string  `json:"name" form:"name"`
}

type AppCUSt struct {
	Name       *string `json:"name" db:"name"`
	PermPrefix *string `json:"perm_prefix" db:"perm_prefix"`
	PermUrl    *string `json:"perm_url" db:"perm_url"`
}
