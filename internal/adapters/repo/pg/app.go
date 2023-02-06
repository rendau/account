package pg

import (
	"context"
	"errors"

	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/dop/adapters/db"
	"github.com/rendau/dop/dopErrs"
)

func (d *St) AppGet(ctx context.Context, id int64) (*entities.AppSt, error) {
	result := &entities.AppSt{}

	err := d.HfGet(ctx, db.RDBGetOptions{
		Dst:    result,
		Tables: []string{"app"},
		Conds:  []string{"id = ${id}"},
		Args:   map[string]any{"id": id},
	})
	if errors.Is(err, dopErrs.NoRows) {
		result = nil
		err = nil
	}

	return result, err
}

func (d *St) AppList(ctx context.Context, pars *entities.AppListParsSt) ([]*entities.AppSt, int64, error) {
	conds := make([]string, 0)
	args := map[string]any{}

	// filter
	if pars.Ids != nil {
		conds = append(conds, `t.id in (select * from unnest(${ids} :: bigint[]))`)
		args["ids"] = *pars.Ids
	}
	if pars.Name != nil {
		conds = append(conds, `t.name = ${name}`)
		args["name"] = *pars.Name
	}

	result := make([]*entities.AppSt, 0, 100)

	tCount, err := d.HfList(ctx, db.RDBListOptions{
		Dst:    &result,
		Tables: []string{`app t`},
		LPars:  pars.ListParams,
		Conds:  conds,
		Args:   args,
		AllowedSorts: map[string]string{
			"default": "t.id",
		},
	})

	return result, tCount, err
}

func (d *St) AppIdExists(ctx context.Context, id int64) (bool, error) {
	var cnt int

	err := d.DbQueryRow(ctx, `
        select count(*)
        from app
        where id = $1
    `, id).Scan(&cnt)

	return cnt > 0, err
}

func (d *St) AppCreate(ctx context.Context, obj *entities.AppCUSt) (int64, error) {
	var result int64

	err := d.HfCreate(ctx, db.RDBCreateOptions{
		Table:  "app",
		Obj:    obj,
		RetCol: "id",
		RetV:   &result,
	})

	return result, err
}

func (d *St) AppUpdate(ctx context.Context, id int64, obj *entities.AppCUSt) error {
	return d.HfUpdate(ctx, db.RDBUpdateOptions{
		Table: "app",
		Obj:   obj,
		Conds: []string{"id = ${cond_id}"},
		Args:  map[string]any{"cond_id": id},
	})
}

func (d *St) AppDelete(ctx context.Context, id int64) error {
	return d.HfDelete(ctx, db.RDBDeleteOptions{
		Table: "app",
		Conds: []string{"id = ${cond_id}"},
		Args:  map[string]any{"cond_id": id},
	})
}
