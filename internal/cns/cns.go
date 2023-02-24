package cns

import "time"

const (
	AppName = "Account"

	MaxPageSize = 1000
)

var (
	AppTimeLocation = time.FixedZone("AST", 21600) // +0600
	True            = true
	False           = false
)

const (
	SFDUsrAva = "usr_avatar"
)

const (
	RoleCodeSuperAdmin = "acc:super_admin"
	RoleCodeAdmin      = "acc:admin"
)

const (
	PermAll   = "acc:*"
	PermMApp  = "acc:m_app"
	PermMPerm = "acc:m_perm"
	PermMRole = "acc:m_role"
	PermMUsr  = "acc:m_usr"
)

const (
	NfTypeRefreshProfile = "refresh-profile"
)
