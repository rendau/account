package usecases

import (
	"context"

	"github.com/rendau/account/internal/domain/entities"
)

func (u *St) DicGet(ctx context.Context) (*entities.DicDataSt, error) {
	return u.cr.Dic.Get(ctx)
}
