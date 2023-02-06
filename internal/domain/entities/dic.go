package entities

type DicDataSt struct {
	Roles []*RoleListSt `json:"roles"`
	Perms []*PermSt     `json:"perms"`
}
