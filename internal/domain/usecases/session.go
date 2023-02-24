package usecases

import (
	"context"

	"github.com/rendau/account/internal/cns"
	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/dop/dopErrs"
)

func (u *St) SessionGetFromToken(token string) *entities.Session {
	return u.cr.Session.GetFromToken(token)
}

func (u *St) SessionRequireAuth(ses *entities.Session) error {
	if ses.Id == 0 {
		return dopErrs.NotAuthorized
	}

	return nil
}

func (u *St) SessionRequirePerm(ses *entities.Session, strict bool, perm string) error {
	err := u.SessionRequireAuth(ses)
	if err != nil {
		return err
	}

	for _, sPerm := range ses.Perms {
		if !strict && sPerm == cns.PermAll {
			return nil
		}
		if perm == sPerm {
			return nil
		}
	}

	return dopErrs.PermissionDenied
}

func (u *St) SessionRequireSAdmin(ses *entities.Session) error {
	if !u.SessionIsSAdmin(ses) {
		return dopErrs.PermissionDenied
	}

	return nil
}

func (u *St) SessionIsSAdmin(ses *entities.Session) bool {
	for _, sPerm := range ses.Perms {
		if sPerm == cns.PermAll {
			return true
		}
	}

	return false
}

func (u *St) SessionSetToContext(ctx context.Context, ses *entities.Session) context.Context {
	return u.cr.Session.SetToContext(ctx, ses)
}

func (u *St) SessionSetToContextByToken(ctx context.Context, token string) context.Context {
	return u.cr.Session.SetToContext(ctx, u.SessionGetFromToken(token))
}

func (u *St) SessionGetFromContext(ctx context.Context) *entities.Session {
	return u.cr.Session.GetFromContext(ctx)
}
