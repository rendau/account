package core

import (
	"context"

	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/account/internal/domain/errs"
	"github.com/rendau/dop/dopErrs"
)

type Perm struct {
	r *St
}

func NewPerm(r *St) *Perm {
	return &Perm{r: r}
}

func (c *Perm) ValidateCU(ctx context.Context, obj *entities.PermCUSt, id int64) error {
	forCreate := id == 0

	// Code
	if forCreate && obj.Code == nil {
		return errs.CodeRequired
	}
	if obj.Code != nil {
		if *obj.Code == "" {
			return errs.CodeRequired
		}
	}

	// AppId
	if forCreate && obj.AppId == nil {
		return errs.ApplicationRequired
	}
	if obj.AppId != nil {
		exists, err := c.r.App.IdExists(ctx, *obj.AppId)
		if err != nil {
			return err
		}
		if !exists {
			return dopErrs.ObjectNotFound
		}
	}

	return nil
}

func (c *Perm) List(ctx context.Context, pars *entities.PermListParsSt) ([]*entities.PermSt, error) {
	items, err := c.r.repo.PermList(ctx, pars)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *Perm) Get(ctx context.Context, id int64, errNE bool) (*entities.PermSt, error) {
	result, err := c.r.repo.PermGet(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		if errNE {
			return nil, dopErrs.ObjectNotFound
		}
		return nil, nil
	}

	return result, nil
}

func (c *Perm) IdExists(ctx context.Context, id int64) (bool, error) {
	return c.r.repo.PermIdExists(ctx, id)
}

func (c *Perm) Create(ctx context.Context, obj *entities.PermCUSt) (int64, error) {
	var err error

	err = c.ValidateCU(ctx, obj, 0)
	if err != nil {
		return 0, err
	}

	// create
	result, err := c.r.repo.PermCreate(ctx, obj)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (c *Perm) Update(ctx context.Context, id int64, obj *entities.PermCUSt) error {
	var err error

	// check if system
	perm, err := c.Get(ctx, id, true)
	if err != nil {
		return err
	}
	if perm.IsSystem {
		return dopErrs.PermissionDenied
	}

	err = c.ValidateCU(ctx, obj, id)
	if err != nil {
		return err
	}

	err = c.r.repo.PermUpdate(ctx, id, obj)
	if err != nil {
		return err
	}

	return nil
}

func (c *Perm) Delete(ctx context.Context, id int64) error {
	// check if system
	perm, err := c.Get(ctx, id, true)
	if err != nil {
		return err
	}
	if perm.IsSystem {
		return dopErrs.PermissionDenied
	}

	return c.r.repo.PermDelete(ctx, id)
}
