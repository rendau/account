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
	UsrGetRoleIds(ctx context.Context, id int64) ([]string, error)
	UsrGetPermIds(ctx context.Context, id int64) ([]string, error)
	UsrGetPhone(ctx context.Context, id int64) (string, error)
	UsrGetIdForPhone(ctx context.Context, phone string) (int64, error)
	UsrCreate(ctx context.Context, obj *entities.UsrCUSt) (int64, error)
	UsrUpdate(ctx context.Context, id int64, obj *entities.UsrCUSt) error
	UsrDelete(ctx context.Context, id int64) error
	UsrFilterUnusedFiles(ctx context.Context, src []string) ([]string, error)

	// role
	RoleGet(ctx context.Context, id string) (*entities.RoleSt, error)
	RoleList(ctx context.Context) ([]*entities.RoleListSt, error)
	RoleCreate(ctx context.Context, obj *entities.RoleCUSt) (string, error)
	RoleIdExists(ctx context.Context, id string) (bool, error)
	RoleUpdate(ctx context.Context, id string, obj *entities.RoleCUSt) error
	RoleDelete(ctx context.Context, id string) error

	// perm
	PermGet(ctx context.Context, id string) (*entities.PermSt, error)
	PermList(ctx context.Context, pars *entities.PermListParsSt) ([]*entities.PermListSt, error)
	PermIdExists(ctx context.Context, id string) (bool, error)
	PermCreate(ctx context.Context, obj *entities.PermCUSt) (string, error)
	PermUpdate(ctx context.Context, id string, obj *entities.PermCUSt) error
	PermDelete(ctx context.Context, id string) error
}
