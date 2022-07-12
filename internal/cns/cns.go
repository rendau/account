package cns

import "time"

const (
	AppName = "Account"

	MaxPageSize = 1000
)

var (
	AppTimeLocation = time.FixedZone("AST", 21600) // +0600
)

const (
	SFDUsrAva = "usr_avatar"
)

const (
	RoleAdmin = "admin"
)

const (
	PermAll   = "*"
	PermMPerm = "m_perm"
	PermMRole = "m_role"
	PermMUsr  = "m_usr"
)

const (
	NfTypeRefreshProfile = "refresh-profile"
)
