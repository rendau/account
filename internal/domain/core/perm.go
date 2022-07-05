package core

import (
	"context"

	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/dop/dopErrs"
)

type Perm struct {
	r *St
}

func NewPerm(r *St) *Perm {
	return &Perm{r: r}
}

func (c *Perm) ValidateCU(ctx context.Context, obj *entities.PermCUSt, id string) error {
	// forCreate := id == ""

	return nil
}

func (c *Perm) List(ctx context.Context, pars *entities.PermListParsSt) ([]*entities.PermListSt, error) {
	items, err := c.r.repo.PermList(ctx, pars)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *Perm) Get(ctx context.Context, id string, errNE bool) (*entities.PermSt, error) {
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

func (c *Perm) IdExists(ctx context.Context, id string) (bool, error) {
	return c.r.repo.PermIdExists(ctx, id)
}

func (c *Perm) Create(ctx context.Context, obj *entities.PermCUSt) (string, error) {
	var err error

	err = c.ValidateCU(ctx, obj, "")
	if err != nil {
		return "", err
	}

	// create
	result, err := c.r.repo.PermCreate(ctx, obj)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (c *Perm) Update(ctx context.Context, id string, obj *entities.PermCUSt) error {
	var err error

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

func (c *Perm) Delete(ctx context.Context, id string) error {
	return c.r.repo.PermDelete(ctx, id)
}
