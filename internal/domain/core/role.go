package core

import (
	"context"

	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/dop/dopErrs"
)

type Role struct {
	r *St
}

func NewRole(r *St) *Role {
	return &Role{r: r}
}

func (c *Role) ValidateCU(ctx context.Context, obj *entities.RoleCUSt, id string) error {
	// forCreate := id == ""

	return nil
}

func (c *Role) List(ctx context.Context) ([]*entities.RoleListSt, error) {
	items, err := c.r.repo.RoleList(ctx)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *Role) Get(ctx context.Context, id string, errNE bool) (*entities.RoleSt, error) {
	result, err := c.r.repo.RoleGet(ctx, id)
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

func (c *Role) IdExists(ctx context.Context, id string) (bool, error) {
	return c.r.repo.RoleIdExists(ctx, id)
}

func (c *Role) Create(ctx context.Context, obj *entities.RoleCUSt) (string, error) {
	err := c.ValidateCU(ctx, obj, "")
	if err != nil {
		return "", err
	}

	return c.r.repo.RoleCreate(ctx, obj)
}

func (c *Role) Update(ctx context.Context, id string, obj *entities.RoleCUSt) error {
	var err error

	err = c.ValidateCU(ctx, obj, id)
	if err != nil {
		return err
	}

	return c.r.repo.RoleUpdate(ctx, id, obj)
}

func (c *Role) Delete(ctx context.Context, id string) error {
	return c.r.repo.RoleDelete(ctx, id)
}
