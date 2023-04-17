package tests

import (
	"github.com/rendau/account/internal/adapters/repo/pg"
	"github.com/rendau/account/internal/cns"
	"github.com/rendau/account/internal/domain/core"
	"github.com/rendau/account/internal/domain/usecases"
	"github.com/rendau/dop/adapters/cache/mem"
	dopDbPg "github.com/rendau/dop/adapters/db/pg"
	dopJwtMock "github.com/rendau/dop/adapters/jwt/mock"
	dopLoggerZap "github.com/rendau/dop/adapters/logger/zap"
	dopSmsMock "github.com/rendau/dop/adapters/sms/mock"
)

var (
	app = struct {
		lg    *dopLoggerZap.St
		cache *mem.St
		db    *dopDbPg.St
		repo  *pg.St
		jwts  *dopJwtMock.St
		sms   *dopSmsMock.St
		core  *core.St
		ucs   *usecases.St
	}{}

	perm1Id = "perm1"
	perm2Id = "perm2"
	perm3Id = "perm3"

	role1Id    = "role1"
	role1Name  = "Role1"
	role1Perms = []string{perm1Id, perm2Id}

	role2Id    = "role2"
	role2Name  = "Role2"
	role2Perms = []string{perm1Id, perm3Id}

	admId    int64
	admName  = "Admin"
	admPhone = "70000000001"
	admRoles = []string{cns.RoleCodeAdmin}

	usr1Id    int64
	usr1Name  = "Usr1"
	usr1Phone = "75550000001"
	usr1Roles = []string{role1Id}
)
