package pg

import (
	"context"
	"errors"

	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/dop/adapters/db"
	"github.com/rendau/dop/dopErrs"
)

func (d *St) PermGet(ctx context.Context, id int64) (*entities.PermSt, error) {
	result := &entities.PermSt{}

	err := d.HfGet(ctx, db.RDBGetOptions{
		Dst:    result,
		Tables: []string{"perm"},
		Conds:  []string{"id = ${id}"},
		Args:   map[string]any{"id": id},
	})
	if err != nil {
		if errors.Is(err, dopErrs.NoRows) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (d *St) PermList(ctx context.Context, pars *entities.PermListParsSt) ([]*entities.PermSt, error) {
	conds := make([]string, 0)
	args := map[string]any{}

	// filter
	if pars.Ids != nil {
		conds = append(conds, `id in (select * from unnest(${ids} :: bigint[]))`)
		args["ids"] = *pars.Ids
	}
	if pars.AppId != nil {
		conds = append(conds, `app_id = ${app_id}`)
		args["app_id"] = *pars.AppId
	}
	if pars.Code != nil {
		conds = append(conds, `code = ${code}`)
		args["code"] = *pars.Code
	}
	if pars.IsSystem != nil {
		if *pars.IsSystem {
			conds = append(conds, `is_system is true`)
		} else {
			conds = append(conds, `is_system is false`)
		}
	}

	result := make([]*entities.PermSt, 0)

	_, err := d.HfList(ctx, db.RDBListOptions{
		Dst:          &result,
		Tables:       []string{`perm`},
		LPars:        pars.ListParams,
		Conds:        conds,
		Args:         args,
		AllowedSorts: map[string]string{"default": "is_system, app_id, code"},
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *St) PermIdExists(ctx context.Context, id int64) (bool, error) {
	var cnt int

	err := d.DbQueryRow(ctx, `
		select count(*)
		from perm
		where id = $1
	`, id).Scan(&cnt)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

func (d *St) PermCreate(ctx context.Context, obj *entities.PermCUSt) (int64, error) {
	var result int64

	err := d.HfCreate(ctx, db.RDBCreateOptions{
		Table:  "perm",
		Obj:    obj,
		RetCol: "id",
		RetV:   &result,
	})
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (d *St) PermUpdate(ctx context.Context, id int64, obj *entities.PermCUSt) error {
	return d.HfUpdate(ctx, db.RDBUpdateOptions{
		Table: "perm",
		Obj:   obj,
		Conds: []string{"id = ${cond_id}"},
		Args:  map[string]any{"cond_id": id},
	})
}

func (d *St) PermDelete(ctx context.Context, id int64) error {
	return d.HfDelete(ctx, db.RDBDeleteOptions{
		Table: "perm",
		Conds: []string{"id = ${cond_id}"},
		Args:  map[string]any{"cond_id": id},
	})
}
