package usecases

import (
	"context"

	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/dop/dopErrs"
)

func (u *St) ProfileSendPhoneValidatingCode(ctx context.Context,
	phone string, errNE bool) error {
	return u.cr.Usr.SendPhoneValidatingCode(ctx, phone, errNE)
}

func (u *St) ProfileAuth(ctx context.Context,
	obj *entities.PhoneAndSmsCodeSt) (string, string, error) {
	var err error

	var accessToken, refreshToken string

	err = u.db.TransactionFn(ctx, func(ctx context.Context) error {
		accessToken, refreshToken, err = u.cr.Usr.Auth(ctx, obj)
		return err
	})

	return accessToken, refreshToken, err
}

func (u *St) ProfileAuthByRefreshToken(ctx context.Context,
	token string) (string, error) {
	return u.cr.Usr.AuthByRefreshToken(ctx, token)
}

func (u *St) ProfileReg(ctx context.Context,
	obj *entities.UsrRegReqSt) (string, string, error) {
	var err error

	var accessToken, refreshToken string

	err = u.db.TransactionFn(ctx, func(ctx context.Context) error {
		accessToken, refreshToken, err = u.cr.Usr.Reg(ctx, obj)
		return err
	})

	return accessToken, refreshToken, err
}

func (u *St) ProfileLogout(ctx context.Context) error {
	ses := u.SessionGetFromContext(ctx)

	if ses.Id == 0 {
		return nil
	}

	return u.db.TransactionFn(ctx, func(ctx context.Context) error {
		return u.cr.Usr.Logout(ctx, ses.Id)
	})
}

func (u *St) ProfileGet(ctx context.Context) (*entities.UsrProfileSt, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return nil, err
	}

	result, err := u.cr.Usr.GetProfile(ctx, ses.Id)
	if err != nil {
		if err == dopErrs.ObjectNotFound {
			err = dopErrs.NotAuthorized
		}
	}

	return result, err
}

func (u *St) ProfileUpdate(ctx context.Context,
	obj *entities.UsrCUSt) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return err
	}

	// restrict
	obj.Roles = nil
	obj.Phone = nil

	return u.db.TransactionFn(ctx, func(ctx context.Context) error {
		return u.cr.Usr.Update(ctx, ses.Id, obj)
	})
}

func (u *St) ProfileChangePhone(ctx context.Context,
	obj *entities.PhoneAndSmsCodeSt) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return err
	}

	return u.db.TransactionFn(ctx, func(ctx context.Context) error {
		return u.cr.Usr.ChangePhone(ctx, ses.Id, obj)
	})
}

func (u *St) ProfileDelete(ctx context.Context) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return err
	}

	return u.db.TransactionFn(ctx, func(ctx context.Context) error {
		return u.cr.Usr.Delete(ctx, ses.Id)
	})
}
