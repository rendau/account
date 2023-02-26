package entities

import (
	"time"

	"github.com/rendau/account/internal/cns"
	"github.com/rendau/dop/dopTypes"
)

type UsrSt struct {
	Id                     int64     `json:"id" db:"id"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	Phone                  string    `json:"phone" db:"phone"`
	Ava                    string    `json:"ava" db:"ava"`
	Name                   string    `json:"name" db:"name"`
	AccessTokenDurSeconds  int64     `json:"access_token_dur_seconds" db:"access_token_dur_seconds"`
	RefreshTokenDurSeconds int64     `json:"refresh_token_dur_seconds" db:"refresh_token_dur_seconds"`

	RoleIds []int64 `json:"role_ids" db:"role_ids"`
	PermIds []int64 `json:"perm_ids" db:"perm_ids"`

	Roles []*RoleListSt `json:"roles,omitempty"`
	Perms []*PermSt     `json:"perms,omitempty"`
}

type UsrGetParsSt struct {
	Id    *int64  `json:"id" form:"id"`
	Phone *string `json:"phone" form:"phone"`

	WithRoles bool `json:"with_roles" form:"with_roles"`
	WithPerms bool `json:"with_perms" form:"with_perms"`
}

type UsrProfileSt struct {
	UsrSt

	RoleCodes []string `json:"role_codes" db:"role_codes"`
	PermCodes []string `json:"perm_codes" db:"perm_codes"`
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
	Phone                  *string `json:"phone" db:"phone"`
	Name                   *string `json:"name" db:"name"`
	Ava                    *string `json:"ava" db:"ava"`
	AccessTokenDurSeconds  *int64  `json:"access_token_dur_seconds" db:"access_token_dur_seconds"`
	RefreshTokenDurSeconds *int64  `json:"refresh_token_dur_seconds" db:"refresh_token_dur_seconds"`

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

type GetNewTokenRepSt struct {
	Token string `json:"token"`
}

func (u *UsrSt) GetRoleCodes() []string {
	res := make([]string, len(u.Roles))
	for i, r := range u.Roles {
		res[i] = r.Code
	}
	return res
}

func (u *UsrSt) HasSAdminRole() bool {
	for _, r := range u.Roles {
		if r.Code == cns.RoleCodeSuperAdmin {
			return true
		}
	}
	return false
}

func (u *UsrSt) GetPermCodes() []string {
	res := make([]string, len(u.Perms))
	for i, p := range u.Perms {
		res[i] = p.Code
	}
	return res
}
