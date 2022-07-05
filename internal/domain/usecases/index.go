package usecases

import (
	"github.com/rendau/account/internal/domain/core"
	"github.com/rendau/dop/adapters/db"
	"github.com/rendau/dop/adapters/logger"
)

type St struct {
	lg logger.Lite

	db db.RDBContextTransaction
	cr *core.St
}

func New(
	lg logger.Lite,
	db db.RDBContextTransaction,
) *St {
	return &St{
		lg: lg,
		db: db,
	}
}

func (u *St) SetCore(core *core.St) {
	u.cr = core
}
