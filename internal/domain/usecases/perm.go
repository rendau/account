package usecases

import (
	"context"

	"github.com/rendau/account/internal/cns"
	"github.com/rendau/account/internal/domain/entities"
)

func (u *St) PermList(ctx context.Context,
	pars *entities.PermListParsSt) ([]*entities.PermSt, error) {
	// var err error
	//
	// ses := u.SessionGetFromContext(ctx)
	//
	// if err = u.SessionRequireAuth(ses); err != nil {
	// 	return nil, err
	// }

	return u.cr.Perm.List(ctx, pars)
}

func (u *St) PermGet(ctx context.Context, id string) (*entities.PermSt, error) {
	// var err error
	//
	// ses := u.SessionGetFromContext(ctx)
	//
	// if err = u.SessionRequireAuth(ses); err != nil {
	// 	return nil, err
	// }

	return u.cr.Perm.Get(ctx, id, true)
}

func (u *St) PermCreate(ctx context.Context,
	obj *entities.PermCUSt) (string, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequirePerm(ses, false, cns.PermMPerm); err != nil {
		return "", err
	}

	var result string

	err = u.db.TransactionFn(ctx, func(ctx context.Context) error {
		result, err = u.cr.Perm.Create(ctx, obj)
		return err
	})

	return result, err
}

func (u *St) PermUpdate(ctx context.Context,
	id string, obj *entities.PermCUSt) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequirePerm(ses, false, cns.PermMPerm); err != nil {
		return err
	}

	return u.db.TransactionFn(ctx, func(ctx context.Context) error {
		return u.cr.Perm.Update(ctx, id, obj)
	})
}

func (u *St) PermDelete(ctx context.Context,
	id string) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequirePerm(ses, false, cns.PermMPerm); err != nil {
		return err
	}

	return u.db.TransactionFn(ctx, func(ctx context.Context) error {
		return u.cr.Perm.Delete(ctx, id)
	})
}
