package pg

import (
	"context"
	"errors"

	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/dop/adapters/db"
	"github.com/rendau/dop/dopErrs"
)

func (d *St) RoleGet(ctx context.Context, id int64) (*entities.RoleSt, error) {
	result := &entities.RoleSt{}

	err := d.HfGet(ctx, db.RDBGetOptions{
		Dst:    result,
		Tables: []string{`"role"`},
		ColExprs: map[string]string{
			"perm_ids": `(
				select array_agg(distinct perm_id)
				from role_perm
				where role_id = id
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

	if result.PermIds == nil {
		result.PermIds = []int64{}
	}

	return result, nil
}

func (d *St) RoleList(ctx context.Context, pars *entities.RoleListParsSt) ([]*entities.RoleListSt, error) {
	conds := make([]string, 0)
	args := map[string]any{}

	// filter
	if pars.Ids != nil {
		conds = append(conds, `id in (select * from unnest(${ids} :: bigint[]))`)
		args["ids"] = *pars.Ids
	}

	result := make([]*entities.RoleListSt, 0)

	_, err := d.HfList(ctx, db.RDBListOptions{
		Dst:          &result,
		Tables:       []string{`"role"`},
		LPars:        pars.ListParams,
		Conds:        conds,
		Args:         args,
		AllowedSorts: map[string]string{"default": "is_system, code"},
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *St) RoleIdExists(ctx context.Context, id int64) (bool, error) {
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

func (d *St) RoleCreate(ctx context.Context, obj *entities.RoleCUSt) (int64, error) {
	var result int64

	err := d.HfCreate(ctx, db.RDBCreateOptions{
		Table:  `"role"`,
		Obj:    obj,
		RetCol: "id",
		RetV:   &result,
	})
	if err != nil {
		return 0, err
	}

	if len(obj.PermIds) > 0 {
		err = d.DbExec(ctx, `
			insert into role_perm(role_id, perm_id)
			select distinct $1::bigint, x.id
			from unnest($2::bigint[]) x(id)
		`, result, obj.PermIds)
		if err != nil {
			return 0, err
		}
	}

	return result, nil
}

func (d *St) RoleUpdate(ctx context.Context, id int64, obj *entities.RoleCUSt) error {
	err := d.HfUpdate(ctx, db.RDBUpdateOptions{
		Table: `"role"`,
		Obj:   obj,
		Conds: []string{"id = ${cond_id}", "not is_system"},
		Args:  map[string]any{"cond_id": id},
	})
	if err != nil {
		return err
	}

	if obj.PermIds != nil {
		err = d.DbExec(ctx, `
			with q0 as (
				select distinct id from unnest($2::bigint[]) x(id)
			), d as (
				delete from role_perm where role_id = $1 and perm_id not in (select id from q0)
			)
			insert into role_perm(role_id, perm_id)
			select $1, q0.id
			from q0
				left join role_perm rp on rp.role_id = $1 and rp.perm_id = q0.id
			where rp.role_id is null
		`, id, obj.PermIds)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *St) RoleDelete(ctx context.Context, id int64) error {
	return d.HfDelete(ctx, db.RDBDeleteOptions{
		Table: `"role"`,
		Conds: []string{"id = ${cond_id}", "not is_system"},
		Args:  map[string]any{"cond_id": id},
	})
}
