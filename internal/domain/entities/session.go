package entities

type Session struct {
	Sub   string   `json:"sub"`
	Id    int64    `json:"id"`
	Roles []string `json:"roles"`
	Perms []string `json:"perms"`
}
