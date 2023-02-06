package repo

import (
	"context"

	"github.com/rendau/account/internal/domain/entities"
)

type Repo interface {
	// usr
	UsrList(ctx context.Context, pars *entities.UsrListParsSt) ([]*entities.UsrListSt, int64, error)
	UsrGet(ctx context.Context, pars *entities.UsrGetParsSt) (*entities.UsrSt, error)
	UsrIdExists(ctx context.Context, id int64) (bool, error)
	UsrIdsExists(ctx context.Context, ids []int64) (bool, error)
	UsrPhoneExists(ctx context.Context, phone string, excludeId int64) (bool, error)
	UsrGetToken(ctx context.Context, id int64) (string, error)
	UsrSetToken(ctx context.Context, id int64, token string) error
	UsrGetRoleIds(ctx context.Context, id int64) ([]int64, error)
	UsrGetPermIds(ctx context.Context, id int64) ([]int64, error)
	UsrGetPhone(ctx context.Context, id int64) (string, error)
	UsrGetIdForPhone(ctx context.Context, phone string) (int64, error)
	UsrCreate(ctx context.Context, obj *entities.UsrCUSt) (int64, error)
	UsrUpdate(ctx context.Context, id int64, obj *entities.UsrCUSt) error
	UsrDelete(ctx context.Context, id int64) error
	UsrFilterUnusedFiles(ctx context.Context, src []string) ([]string, error)

	// role
	RoleGet(ctx context.Context, id int64) (*entities.RoleSt, error)
	RoleList(ctx context.Context, pars *entities.RoleListParsSt) ([]*entities.RoleListSt, error)
	RoleCreate(ctx context.Context, obj *entities.RoleCUSt) (string, error)
	RoleIdExists(ctx context.Context, id int64) (bool, error)
	RoleUpdate(ctx context.Context, id int64, obj *entities.RoleCUSt) error
	RoleDelete(ctx context.Context, id int64) error

	// perm
	PermGet(ctx context.Context, id int64) (*entities.PermSt, error)
	PermList(ctx context.Context, pars *entities.PermListParsSt) ([]*entities.PermSt, error)
	PermIdExists(ctx context.Context, id int64) (bool, error)
	PermCreate(ctx context.Context, obj *entities.PermCUSt) (string, error)
	PermUpdate(ctx context.Context, id int64, obj *entities.PermCUSt) error
	PermDelete(ctx context.Context, id int64) error

	// app
	AppGet(ctx context.Context, id int64) (*entities.AppSt, error)
	AppList(ctx context.Context, pars *entities.AppListParsSt) ([]*entities.AppSt, int64, error)
	AppIdExists(ctx context.Context, id int64) (bool, error)
	AppCreate(ctx context.Context, obj *entities.AppCUSt) (int64, error)
	AppUpdate(ctx context.Context, id int64, obj *entities.AppCUSt) error
	AppDelete(ctx context.Context, id int64) error
}
