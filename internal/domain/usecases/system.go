package usecases

import (
	"github.com/rendau/account/internal/domain/entities"
)

func (u *St) SystemFilterUnusedFiles(filePaths []string) []string {
	return u.cr.System.FilterUnusedFiles(filePaths)
}

func (u *St) SystemCronTick5m() {
	u.cr.System.CronTick5m()
}

func (u *St) SystemCronTick15m() {
	u.cr.System.CronTick15m()
}

func (u *St) SystemCronTick30m() {
	u.cr.System.CronTick30m()
}

func (u *St) GetPerms() *entities.SystemGetPermsRepSt {
	return u.cr.System.GetPerms()
}
