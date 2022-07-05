package usecases

import (
	"context"

	"github.com/rendau/account/internal/cns"
	"github.com/rendau/account/internal/domain/entities"
)

func (u *St) RoleList(ctx context.Context) ([]*entities.RoleListSt, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return nil, err
	}

	return u.cr.Role.List(ctx)
}

func (u *St) RoleGet(ctx context.Context, id string) (*entities.RoleSt, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return nil, err
	}

	return u.cr.Role.Get(ctx, id, true)
}

func (u *St) RoleCreate(ctx context.Context,
	obj *entities.RoleCUSt) (string, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequirePerm(ses, false, cns.PermMRole); err != nil {
		return "", err
	}

	var result string

	err = u.db.TransactionFn(ctx, func(ctx context.Context) error {
		result, err = u.cr.Role.Create(ctx, obj)
		return err
	})

	return result, err
}

func (u *St) RoleUpdate(ctx context.Context,
	id string, obj *entities.RoleCUSt) error {
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
	id string) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequirePerm(ses, false, cns.PermMRole); err != nil {
		return err
	}

	return u.db.TransactionFn(ctx, func(ctx context.Context) error {
		return u.cr.Role.Delete(ctx, id)
	})
}
