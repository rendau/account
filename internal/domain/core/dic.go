package core

import (
	"context"

	"github.com/rendau/account/internal/domain/entities"
)

type Dic struct {
	r *St
}

func NewDic(r *St) *Dic {
	return &Dic{r: r}
}

func (c *Dic) Get(ctx context.Context) (*entities.DicDataSt, error) {
	var err error

	data := &entities.DicDataSt{}

	data.Roles, err = c.r.Role.List(ctx, &entities.RoleListParsSt{})
	if err != nil {
		return nil, err
	}

	data.Apps, _, err = c.r.App.List(ctx, &entities.AppListParsSt{})
	if err != nil {
		return nil, err
	}

	return data, nil
}
