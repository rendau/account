package pg

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/dop/adapters/db"
	"github.com/rendau/dop/dopErrs"
)

func (d *St) UsrList(ctx context.Context, pars *entities.UsrListParsSt) ([]*entities.UsrListSt, int64, error) {
	conds := make([]string, 0)
	args := map[string]any{}

	// filter
	if pars.Id != nil {
		conds = append(conds, `id = ${id}`)
		args["id"] = *pars.Id
	}
	if pars.Ids != nil {
		conds = append(conds, `id in (select * from unnest(${ids} :: bigint[]))`)
		args["ids"] = *pars.Ids
	}
	if pars.Search != nil && *pars.Search != "" {
		for wordI, word := range strings.Split(*pars.Search, " ") {
			if word != "" {
				key := "s_word_" + strconv.Itoa(wordI)

				args[key] = word
				conds = append(conds, `(
					phone ilike '%'|| ${`+key+`} ||'%' or
					name ilike '%'|| ${`+key+`} ||'%'
				)`)
			}
		}
	}

	result := make([]*entities.UsrListSt, 0)

	tCount, err := d.HfList(ctx, db.RDBListOptions{
		Dst:    &result,
		Tables: []string{"usr"},
		LPars:  pars.ListParams,
		Conds:  conds,
		Args:   args,
		AllowedSorts: map[string]string{
			"default": `name collate "C", id`,
			"phone":   "phone",
			"name":    "name",
		},
	})
	if err != nil {
		return nil, 0, err
	}

	return result, tCount, nil
}

func (d *St) UsrGet(ctx context.Context, pars *entities.UsrGetParsSt) (*entities.UsrSt, error) {
	conds := make([]string, 0)
	args := map[string]any{}

	// filter
	if pars.Id != nil {
		conds = append(conds, `id = ${id}`)
		args["id"] = *pars.Id
	}
	if pars.Phone != nil {
		conds = append(conds, `phone = ${phone}`)
		args["phone"] = *pars.Phone
	}
	if pars.Token != nil {
		conds = append(conds, `token = ${token}`)
		args["token"] = *pars.Token
	}

	result := &entities.UsrSt{}

	err := d.HfGet(ctx, db.RDBGetOptions{
		Dst:    result,
		Tables: []string{"usr"},
		ColExprs: map[string]string{
			"role_ids": `(
				select array_agg(distinct role_id)
				from usr_role
				where usr_id = usr.id
			)`,
			"perm_ids": `(
				select array_agg(distinct rp.perm_id)
				from usr_role ur
					join role_perm rp on rp.role_id = ur.role_id
				where ur.usr_id = usr.id
			)`,
		},
		Conds: conds,
		Args:  args,
	})
	if err != nil {
		if errors.Is(err, dopErrs.NoRows) {
			return nil, nil
		}
		return nil, err
	}

	if result.RoleIds == nil {
		result.RoleIds = []int64{}
	}
	if result.PermIds == nil {
		result.PermIds = []int64{}
	}

	return result, nil
}

func (d *St) UsrIdExists(ctx context.Context, id int64) (bool, error) {
	var cnt int

	err := d.DbQueryRow(ctx, `
		select count(*)
		from usr
		where id = $1
	`, id).Scan(&cnt)

	return cnt > 0, err
}

func (d *St) UsrIdsExists(ctx context.Context, ids []int64) (bool, error) {
	var cnt int

	err := d.DbQueryRow(ctx, `
		select count(*)
		from unnest($1::bigint[]) x(id)
			left join usr u on u.id = x.id
		where u.id is null
	`, ids).Scan(&cnt)

	return cnt == 0, err
}

func (d *St) UsrPhoneExists(ctx context.Context, phone string, excludeId int64) (bool, error) {
	var cnt int

	err := d.DbQueryRow(ctx, `
		select count(*)
		from usr
		where phone = $1
			and id != $2
	`, phone, excludeId).Scan(&cnt)

	return cnt > 0, err
}

func (d *St) UsrGetToken(ctx context.Context, id int64) (string, error) {
	var token string

	err := d.DbQueryRow(ctx, `
		select token
		from usr
		where id = $1
	`, id).Scan(&token)

	return token, err
}

func (d *St) UsrSetToken(ctx context.Context, id int64, token string) error {
	return d.DbExec(ctx, `update usr set token = $2 where id = $1`, id, token)
}

func (d *St) UsrGetRoleIds(ctx context.Context, id int64) ([]int64, error) {
	result := make([]int64, 0)

	err := d.DbQueryRow(ctx, `
		select array_agg(distinct role_id)
		from usr_role
		where usr_id = $1
	`, id).Scan(&result)

	return result, err
}

func (d *St) UsrGetPermIds(ctx context.Context, id int64) ([]int64, error) {
	result := make([]int64, 0)

	err := d.DbQueryRow(ctx, `
		select array_agg(distinct rp.perm_id)
		from usr_role ur
			join role_perm rp on rp.role_id = ur.role_id
		where ur.usr_id = $1
	`, id).Scan(&result)

	return result, err
}

func (d *St) UsrGetPhone(ctx context.Context, id int64) (string, error) {
	var result string

	err := d.DbQueryRow(ctx, `
		select phone
		from usr
		where id = $1
	`, id).Scan(&result)

	return result, err
}

func (d *St) UsrGetIdForPhone(ctx context.Context, phone string) (int64, error) {
	var result int64

	err := d.DbQueryRow(ctx, `
		select id
		from usr
		where phone = $1
	`, phone).Scan(&result)
	if err != nil {
		if errors.Is(err, dopErrs.NoRows) {
			return 0, nil
		}

		return 0, err
	}

	return result, nil
}

func (d *St) UsrCreate(ctx context.Context, obj *entities.UsrCUSt) (int64, error) {
	var result int64

	err := d.HfCreate(ctx, db.RDBCreateOptions{
		Table:  "usr",
		Obj:    obj,
		RetCol: "id",
		RetV:   &result,
	})
	if err != nil {
		return 0, err
	}

	if len(obj.RoleIds) > 0 {
		err = d.DbExec(ctx, `
			insert into usr_role(usr_id, role_id)
			select distinct $1::bigint, x.id
			from unnest($2::bigint[]) x(id)
		`, result, obj.RoleIds)
		if err != nil {
			return 0, err
		}
	}

	return result, nil
}

func (d *St) UsrUpdate(ctx context.Context, id int64, obj *entities.UsrCUSt) error {
	err := d.HfUpdate(ctx, db.RDBUpdateOptions{
		Table: "usr",
		Obj:   obj,
		Conds: []string{"id = ${cond_id}"},
		Args:  map[string]any{"cond_id": id},
	})
	if err != nil {
		return err
	}

	if obj.RoleIds != nil {
		err = d.DbExec(ctx, `
			with q0 as (
				select distinct id from unnest($2::bigint[]) x(id)
			), d as (
				delete from usr_role where usr_id = $1 and role_id not in (select id from q0)
			)
			insert into usr_role(usr_id, role_id)
			select $1, q0.id
			from q0
				left join usr_role ur on ur.usr_id = $1 and ur.role_id = q0.id
			where ur.role_id is null
		`, id, obj.RoleIds)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *St) UsrDelete(ctx context.Context, id int64) error {
	return d.HfDelete(ctx, db.RDBDeleteOptions{
		Table: "usr",
		Conds: []string{"id = ${cond_id}"},
		Args:  map[string]any{"cond_id": id},
	})
}

func (d *St) UsrFilterUnusedFiles(ctx context.Context, src []string) ([]string, error) {
	rows, err := d.DbQuery(ctx, `
		select x.a
		from unnest($1 :: bigint[]) x(a)
			left join usr y on y.ava = x.a
		where y.ava is null
	`, src)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]string, 0)

	var v string

	for rows.Next() {
		err = rows.Scan(&v)
		if err != nil {
			return nil, err
		}

		result = append(result, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
