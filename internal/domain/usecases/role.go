package usecases

import (
	"context"

	"github.com/rendau/account/internal/cns"
	"github.com/rendau/account/internal/domain/entities"
)

func (u *St) RoleList(ctx context.Context,
	pars *entities.RoleListParsSt) ([]*entities.RoleListSt, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return nil, err
	}

	return u.cr.Role.List(ctx, pars)
}

func (u *St) RoleGet(ctx context.Context,
	id int64) (*entities.RoleSt, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return nil, err
	}

	return u.cr.Role.Get(ctx, id, true)
}

func (u *St) RoleCreate(ctx context.Context,
	obj *entities.RoleCUSt) (int64, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequirePerm(ses, false, cns.PermMRole); err != nil {
		return 0, err
	}

	var result int64

	err = u.db.TransactionFn(ctx, func(ctx context.Context) error {
		result, err = u.cr.Role.Create(ctx, obj)
		return err
	})

	return result, err
}

func (u *St) RoleUpdate(ctx context.Context,
	id int64, obj *entities.RoleCUSt) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequirePerm(ses, false, cns.PermMRole); err != nil {
		return err
	}

	return u.db.TransactionFn(ctx, func(ctx context.Context) error {
		return u.cr.Role.Update(ctx, id, obj)
	})
}

func (u *St) RoleDelete(ctx context.Context,
	id int64) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequirePerm(ses, false, cns.PermMRole); err != nil {
		return err
	}

	return u.db.TransactionFn(ctx, func(ctx context.Context) error {
		return u.cr.Role.Delete(ctx, id)
	})
}
