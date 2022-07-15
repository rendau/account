package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/rendau/account/internal/cns"
	"github.com/rendau/account/internal/domain/entities"
	"github.com/stretchr/testify/require"
)

func resetDb() {
	var err error

	err = app.db.DbExec(context.Background(), `delete from perm where not is_system`)
	if err != nil {
		app.lg.Fatal(err)
	}

	err = app.db.DbExec(context.Background(), `delete from role where not is_system`)
	if err != nil {
		app.lg.Fatal(err)
	}

	truncateTables([]string{
		"usr",
	})

	bgCtx := context.Background()

	perms := []string{
		perm1Id,
		perm2Id,
		perm3Id,
	}
	application := "application1"
	for _, permId := range perms {
		_, err = app.core.Perm.Create(bgCtx, &entities.PermCUSt{
			Id:  &permId,
			App: &application,
			Dsc: &permId,
		})
		if err != nil {
			fmt.Println(permId)
			app.lg.Fatal(err)
		}
	}

	roles := []entities.RoleCUSt{
		{Id: &role1Id, Name: &role1Name, Perms: role1Perms},
		{Id: &role2Id, Name: &role2Name, Perms: role2Perms},
	}
	for _, role := range roles {
		_, err = app.core.Role.Create(bgCtx, &role)
		if err != nil {
			app.lg.Fatal(err)
		}
	}

	usrs := []struct {
		IdPtr *int64
		Name  string
		Phone string
		Roles []string
	}{
		{&admId, admName, admPhone, admRoles},
		{&usr1Id, usr1Name, usr1Phone, usr1Roles},
	}
	for _, usr := range usrs {
		*usr.IdPtr, err = app.core.Usr.Create(bgCtx, &entities.UsrCUSt{
			Roles: usr.Roles,
			Name:  &usr.Name,
			Phone: &usr.Phone,
		})
		if err != nil {
			app.lg.Fatal(err)
		}
	}
}

func prepareDbForNewTest() {
	var err error

	app.cache.Clean()
	app.sms.Clean()

	// truncateTables([]string{})

	err = app.db.DbExec(context.Background(), `
		delete from perm where id not in (select * from unnest($1 :: text[])) and not is_system
	`, []string{perm1Id, perm2Id, perm3Id})
	if err != nil {
		app.lg.Fatal(err)
	}

	err = app.db.DbExec(context.Background(), `
		delete from "role" where id not in (select * from unnest($1 :: text[])) and not is_system
	`, []string{cns.RoleAdmin, role1Id, role2Id})
	if err != nil {
		app.lg.Fatal(err)
	}

	err = app.db.DbExec(context.Background(), `
		delete from usr where id not in (select * from unnest($1 :: bigint[]))
	`, []int64{admId, usr1Id})
	if err != nil {
		app.lg.Fatal(err)
	}
}

func ctxWithSes(t *testing.T, ctx context.Context, usrId int64) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	roleIds, err := app.core.Usr.GetRoleIds(ctx, usrId)
	require.Nil(t, err)

	permIds, err := app.core.Usr.GetPermIds(ctx, usrId)
	require.Nil(t, err)

	return app.ucs.SessionSetToContext(ctx, &entities.Session{
		Id:    usrId,
		Roles: roleIds,
		Perms: permIds,
	})
}

func truncateTables(tables []string) {
	if len(tables) == 0 {
		return
	}

	q := ``
	for _, t := range tables {
		q += ` truncate ` + t + ` restart identity cascade; `
	}
	if q != `` {
		err := app.db.DbExec(context.Background(), `begin; `+q+` commit;`)
		if err != nil {
			app.lg.Fatal(err)
		}
	}
}
