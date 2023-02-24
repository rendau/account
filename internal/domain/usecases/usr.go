package usecases

import (
	"context"

	"github.com/rendau/account/internal/cns"
	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/dop/dopErrs"
	"github.com/rendau/dop/dopTools"
	"github.com/rendau/dop/dopTypes"
)

func (u *St) UsrList(ctx context.Context,
	pars *entities.UsrListParsSt) ([]*entities.UsrListSt, int64, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return nil, 0, err
	}

	if err = dopTools.RequirePageSize(pars.ListParams, cns.MaxPageSize); err != nil {
		return nil, 0, err
	}

	return u.cr.Usr.List(ctx, pars)
}

func (u *St) UsrGet(ctx context.Context,
	id int64) (*entities.UsrSt, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return nil, err
	}

	return u.cr.Usr.Get(ctx, &entities.UsrGetParsSt{Id: &id}, true)
}

func (u *St) UsrCreate(ctx context.Context,
	obj *entities.UsrCUSt) (int64, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequirePerm(ses, false, cns.PermMUsr); err != nil {
		return 0, err
	}

	// check if super admin role is not set
	if !u.SessionIsSAdmin(ses) && len(obj.RoleIds) > 0 {
		items, err := u.cr.Role.List(ctx, &entities.RoleListParsSt{
			ListParams: dopTypes.ListParams{PageSize: 1},
			Ids:        &obj.RoleIds,
			Code:       dopTools.NewPtr(cns.RoleCodeSuperAdmin),
		})
		if err != nil {
			return 0, err
		}
		if len(items) > 0 {
			return 0, dopErrs.PermissionDenied
		}
	}

	var result int64

	err = u.db.TransactionFn(ctx, func(ctx context.Context) error {
		result, err = u.cr.Usr.Create(ctx, obj)
		return err
	})

	return result, err
}

func (u *St) UsrUpdate(ctx context.Context,
	id int64, obj *entities.UsrCUSt) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequirePerm(ses, false, cns.PermMUsr); err != nil {
		return err
	}

	// check if super admin role is not changed
	if !u.SessionIsSAdmin(ses) && obj.RoleIds != nil {
		usrIsSAdmin, err := u.cr.Usr.IsSAdmin(ctx, id)
		if err != nil {
			return err
		}

		items, err := u.cr.Role.List(ctx, &entities.RoleListParsSt{
			ListParams: dopTypes.ListParams{PageSize: 1},
			Ids:        &obj.RoleIds,
			Code:       dopTools.NewPtr(cns.RoleCodeSuperAdmin),
		})
		if err != nil {
			return err
		}
		if len(items) > 0 != usrIsSAdmin {
			return dopErrs.PermissionDenied
		}
	}

	return u.db.TransactionFn(ctx, func(ctx context.Context) error {
		return u.cr.Usr.Update(ctx, id, obj)
	})
}

func (u *St) UsrDelete(ctx context.Context,
	id int64) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequirePerm(ses, false, cns.PermMUsr); err != nil {
		return err
	}

	// check if super admin deleted
	if !u.SessionIsSAdmin(ses) {
		usrIsSAdmin, err := u.cr.Usr.IsSAdmin(ctx, id)
		if err != nil {
			return err
		}
		if usrIsSAdmin {
			return dopErrs.PermissionDenied
		}
	}

	return u.db.TransactionFn(ctx, func(ctx context.Context) error {
		return u.cr.Usr.Delete(ctx, id)
	})
}
