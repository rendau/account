package entities

import (
	"time"

	"github.com/rendau/dop/dopTypes"
)

type UsrSt struct {
	Id        int64     `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Phone     string    `json:"phone" db:"phone"`
	Ava       string    `json:"ava" db:"ava"`
	Name      string    `json:"name" db:"name"`

	RoleIds []int64       `json:"role_ids" db:"role_ids"`
	Roles   []*RoleListSt `json:"roles" db:"roles"`
	PermIds []int64       `json:"perm_ids" db:"perm_ids"`
	Perms   []*PermSt     `json:"perms" db:"perms"`
}

type UsrGetParsSt struct {
	Id    *int64  `json:"id" form:"id"`
	Phone *string `json:"phone" form:"phone"`
	Token *string `json:"token" form:"token"`

	WithRoles bool `json:"with_roles" form:"with_roles"`
	WithPerms bool `json:"with_perms" form:"with_perms"`
}

type UsrProfileSt struct {
	UsrSt
}

type UsrListSt struct {
	Id        int64     `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Phone     string    `json:"phone" db:"phone"`
	Ava       string    `json:"ava" db:"ava"`
	Name      string    `json:"name" db:"name"`
}

type UsrListParsSt struct {
	dopTypes.ListParams

	Id     *int64   `json:"id" form:"id"`
	Ids    *[]int64 `json:"ids" form:"ids"`
	Search *string  `json:"search" form:"search"`
}

type UsrCUSt struct {
	Phone *string `json:"phone" db:"phone"`
	Name  *string `json:"name" db:"name"`
	Ava   *string `json:"ava" db:"ava"`

	RoleIds []int64 `json:"role_ids"`
}

type PhoneAndSmsCodeSt struct {
	Phone   string `json:"phone"`
	SmsCode int    `json:"sms_code"`
}

type UsrRegReqSt struct {
	PhoneAndSmsCodeSt

	Name *string `json:"name"`
	Ava  *string `json:"ava"`
}

type SendPhoneValidatingCodeReqSt struct {
	Phone string `json:"phone"`
	ErrNE bool   `json:"err_ne"`
}

type AuthRepSt struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthByTokenReqSt struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthByTokenRepSt struct {
	AccessToken string `json:"access_token"`
}
