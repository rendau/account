package entities

type RoleSt struct {
	Id       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	IsSystem bool   `json:"is_system" db:"is_system"`

	Perms []string `json:"perms" db:"perms"`
}

type RoleListSt struct {
	Id       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	IsSystem bool   `json:"is_system" db:"is_system"`
}

type RoleCUSt struct {
	Id   *string `json:"id" db:"id"`
	Name *string `json:"name" db:"name"`

	Perms []string `json:"perms"`
}
