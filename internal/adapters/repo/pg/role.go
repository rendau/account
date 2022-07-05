package pg

import (
	"context"
	"errors"

	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/dop/adapters/db"
	"github.com/rendau/dop/dopErrs"
)

func (d *St) RoleGet(ctx context.Context, id string) (*entities.RoleSt, error) {
	result := &entities.RoleSt{}

	err := d.HfGet(ctx, db.RDBGetOptions{
		Dst:    result,
		Tables: []string{`"role"`},
		ColExprs: map[string]string{
			"perms": `(
				select array_agg(distinct perm_id)
				from role_perm
				where role_id = t.id
			)`,
		},
		Conds: []string{"id = ${id}"},
		Args:  map[string]any{"id": id},
	})
	if err != nil {
		if errors.Is(err, dopErrs.NoRows) {
			return nil, nil
		}
		return nil, err
	}

	if result.Perms == nil {
		result.Perms = []string{}
	}

	return result, nil
}

func (d *St) RoleList(ctx context.Context) ([]*entities.RoleListSt, error) {
	result := make([]*entities.RoleListSt, 0)

	_, err := d.HfList(ctx, db.RDBListOptions{
		Dst:          &result,
		Tables:       []string{`"role"`},
		AllowedSorts: map[string]string{"default": "id"},
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *St) RoleIdExists(ctx context.Context, id string) (bool, error) {
	var cnt int

	err := d.DbQueryRow(ctx, `
		select count(*)
		from "role"
		where id = $1
	`, id).Scan(&cnt)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

func (d *St) RoleCreate(ctx context.Context, obj *entities.RoleCUSt) (string, error) {
	var result string

	err := d.HfCreate(ctx, db.RDBCreateOptions{
		Table:  `"role"`,
		Obj:    obj,
		RetCol: "id",
		RetV:   &result,
	})
	if err != nil {
		return "", err
	}

	if len(obj.Perms) > 0 {
		err = d.DbExec(ctx, `
			insert into role_perm(role_id, perm_id)
			select distinct $1, x.id
			from unnest($2::text[]) x(id)
		`, result, obj.Perms)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

func (d *St) RoleUpdate(ctx context.Context, id string, obj *entities.RoleCUSt) error {
	err := d.HfUpdate(ctx, db.RDBUpdateOptions{
		Table: `"role"`,
		Obj:   obj,
		Conds: []string{"id = ${cond_id}", "not is_system"},
		Args:  map[string]any{"cond_id": id},
	})
	if err != nil {
		return err
	}

	if obj.Perms != nil {
		err = d.DbExec(ctx, `
			with q0 as (
				select distinct id from unnest($2::text[]) x(id)
			), d as (
				delete from role_perm where role_id = $1 and perm_id not in (select id from q0)
			)
			insert into role_perm(role_id, perm_id)
			select $1, q0.id
			from q0
				left join role_perm rp on rp.usr_id = $1 and rp.role_id = q0.id
			where rp.role_id is null
		`, id, obj.Perms)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *St) RoleDelete(ctx context.Context, id string) error {
	return d.HfDelete(ctx, db.RDBDeleteOptions{
		Table: `"role"`,
		Conds: []string{"id = ${cond_id}", "not is_system"},
		Args:  map[string]any{"cond_id": id},
	})
}
