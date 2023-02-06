package core

import (
	"context"

	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/dop/dopErrs"
)

type App struct {
	r *St
}

func NewApp(r *St) *App {
	return &App{r: r}
}

func (c *App) ValidateCU(ctx context.Context, obj *entities.AppCUSt, id int64) error {
	// forCreate := id == 0

	return nil
}

func (c *App) List(ctx context.Context, pars *entities.AppListParsSt) ([]*entities.AppSt, int64, error) {
	items, tCount, err := c.r.repo.AppList(ctx, pars)
	if err != nil {
		return nil, 0, err
	}

	return items, tCount, nil
}

func (c *App) Get(ctx context.Context, id int64, errNE bool) (*entities.AppSt, error) {
	result, err := c.r.repo.AppGet(ctx, id)
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

func (c *App) IdExists(ctx context.Context, id int64) (bool, error) {
	return c.r.repo.AppIdExists(ctx, id)
}

func (c *App) Create(ctx context.Context, obj *entities.AppCUSt) (int64, error) {
	var err error

	err = c.ValidateCU(ctx, obj, 0)
	if err != nil {
		return 0, err
	}

	// create
	result, err := c.r.repo.AppCreate(ctx, obj)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (c *App) Update(ctx context.Context, id int64, obj *entities.AppCUSt) error {
	var err error

	err = c.ValidateCU(ctx, obj, id)
	if err != nil {
		return err
	}

	err = c.r.repo.AppUpdate(ctx, id, obj)
	if err != nil {
		return err
	}

	return nil
}

func (c *App) Delete(ctx context.Context, id int64) error {
	return c.r.repo.AppDelete(ctx, id)
}
