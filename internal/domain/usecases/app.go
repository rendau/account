package usecases

import (
	"context"

	"github.com/rendau/account/internal/cns"
	"github.com/rendau/account/internal/domain/entities"
)

func (u *St) AppList(ctx context.Context,
	pars *entities.AppListParsSt) ([]*entities.AppSt, int64, error) {
	// var err error
	//
	// ses := u.SessionGetFromContext(ctx)
	//
	// if err = u.SessionRequireAuth(ses); err != nil {
	// 	return nil, 0, err
	// }

	return u.cr.App.List(ctx, pars)
}

func (u *St) AppGet(ctx context.Context, id int64) (*entities.AppSt, error) {
	// var err error
	//
	// ses := u.SessionGetFromContext(ctx)
	//
	// if err = u.SessionRequireAuth(ses); err != nil {
	// 	return nil, err
	// }

	return u.cr.App.Get(ctx, id, true)
}

func (u *St) AppCreate(ctx context.Context,
	obj *entities.AppCUSt) (int64, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequirePerm(ses, false, cns.PermMApp); err != nil {
		return 0, err
	}

	var result int64

	err = u.db.TransactionFn(ctx, func(ctx context.Context) error {
		result, err = u.cr.App.Create(ctx, obj)
		return err
	})

	return result, err
}

func (u *St) AppUpdate(ctx context.Context,
	id int64, obj *entities.AppCUSt) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequirePerm(ses, false, cns.PermMApp); err != nil {
		return err
	}

	return u.db.TransactionFn(ctx, func(ctx context.Context) error {
		return u.cr.App.Update(ctx, id, obj)
	})
}

func (u *St) AppDelete(ctx context.Context,
	id int64) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequirePerm(ses, false, cns.PermMApp); err != nil {
		return err
	}

	return u.db.TransactionFn(ctx, func(ctx context.Context) error {
		return u.cr.App.Delete(ctx, id)
	})
}
