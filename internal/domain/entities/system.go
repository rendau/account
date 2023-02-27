package entities

type SystemGetPermsRepSt struct {
	Perms []*SystemGetPermsItemSt `json:"perms"`
}

type SystemGetPermsItemSt struct {
	Code  string `json:"code"`
	IsAll bool   `json:"is_all"`
	Dsc   string `json:"dsc"`
}
