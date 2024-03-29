package core

import (
	"context"
	"strings"

	"github.com/rendau/account/internal/cns"
	"github.com/rendau/account/internal/domain/entities"
)

type System struct {
	r *St
}

func NewSystem(r *St) *System {
	return &System{r: r}
}

// FilterUnusedFiles must return files (from 'filePaths') that are not exists in db anymore, and must check only specific directories
func (c *System) FilterUnusedFiles(filePaths []string) []string {
	var err error

	ctx := context.Background()

	modules := []struct {
		filterFn func(context.Context, []string) ([]string, error)
		sfd      string
	}{
		{c.r.repo.UsrFilterUnusedFiles, cns.SFDUsrAva},
	}

	result := make([]string, 0, len(filePaths))

	for _, f := range filePaths {
		for _, module := range modules {
			if strings.HasPrefix(f, module.sfd) {
				result = append(result, f)
				break
			}
		}
	}
	if len(result) == 0 {
		return result
	}

	for _, module := range modules {
		result, err = module.filterFn(ctx, result)
		if err != nil {
			c.r.lg.Errorw("Fail to filter unused files", err)
			return []string{}
		}

		if len(result) == 0 {
			break
		}
	}

	return result
}

func (c *System) CronTick5m() {
	// do something in goroutine
}

func (c *System) CronTick15m() {
	// do something in goroutine
}

func (c *System) CronTick30m() {
	// do something in goroutine
}

func (c *System) GetPerms() *entities.SystemGetPermsRepSt {
	return &entities.SystemGetPermsRepSt{
		Perms: []*entities.SystemGetPermsItemSt{
			{Code: cns.PermAll, IsAll: true, Dsc: "All permissions"},
			{Code: cns.PermMApp, IsAll: false, Dsc: "Modify applications"},
			{Code: cns.PermMPerm, IsAll: false, Dsc: "Modify permissions"},
			{Code: cns.PermMRole, IsAll: false, Dsc: "Modify roles"},
			{Code: cns.PermMUsr, IsAll: false, Dsc: "Modify users"},
		},
	}
}
