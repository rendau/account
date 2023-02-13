package entities

type DicDataSt struct {
	Roles []*RoleListSt `json:"roles"`
	Apps  []*AppSt      `json:"apps"`
}
