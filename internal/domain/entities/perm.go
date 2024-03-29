package entities

import (
	"github.com/rendau/dop/dopTypes"
)

type PermSt struct {
	Id        int64  `json:"id" db:"id"`
	Code      string `json:"code" db:"code"`
	IsAll     bool   `json:"is_all" db:"is_all"`
	AppId     int64  `json:"app_id" db:"app_id"`
	Dsc       string `json:"dsc" db:"dsc"`
	IsSystem  bool   `json:"is_system" db:"is_system"`
	IsFetched bool   `json:"is_fetched" db:"is_fetched"`
}

type PermListParsSt struct {
	dopTypes.ListParams

	Ids       *[]int64 `json:"ids" form:"ids"`
	AppId     *int64   `json:"app_id" form:"app_id"`
	Code      *string  `json:"code" form:"code"`
	IsSystem  *bool    `json:"is_system" form:"is_system"`
	IsFetched *bool    `json:"is_fetched" form:"is_fetched"`
}

type PermCUSt struct {
	Code      *string `json:"code" db:"code"`
	IsAll     *bool   `json:"is_all" db:"is_all"`
	AppId     *int64  `json:"app_id" db:"app_id"`
	Dsc       *string `json:"dsc" db:"dsc"`
	IsFetched *bool   `json:"-" db:"is_fetched"`
}
