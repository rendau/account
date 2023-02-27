package core

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/rendau/account/internal/cns"
	"github.com/rendau/account/internal/domain/entities"
	"github.com/rendau/account/internal/domain/errs"
	"github.com/rendau/dop/adapters/client/httpc"
	"github.com/rendau/dop/adapters/client/httpc/httpclient"
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

	err = c.SyncPerms(ctx, result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (c *App) Update(ctx context.Context, id int64, obj *entities.AppCUSt) error {
	var err error

	// check if system
	perm, err := c.Get(ctx, id, true)
	if err != nil {
		return err
	}
	if perm.IsAccountApp {
		return dopErrs.PermissionDenied
	}

	err = c.ValidateCU(ctx, obj, id)
	if err != nil {
		return err
	}

	err = c.r.repo.AppUpdate(ctx, id, obj)
	if err != nil {
		return err
	}

	err = c.SyncPerms(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *App) Delete(ctx context.Context, id int64) error {
	// check if system
	perm, err := c.Get(ctx, id, true)
	if err != nil {
		return err
	}
	if perm.IsAccountApp {
		return dopErrs.PermissionDenied
	}

	return c.r.repo.AppDelete(ctx, id)
}

func (c *App) FetchPerms(ctx context.Context, uri string) (*entities.SystemGetPermsRepSt, error) {
	return c.fetchPerms(uri)
}

func (c *App) SyncPerms(ctx context.Context, id int64) error {
	app, err := c.Get(ctx, id, true)
	if err != nil {
		return err
	}

	dbPerms, err := c.r.Perm.List(ctx, &entities.PermListParsSt{
		AppId:     &id,
		IsFetched: &cns.True,
	})
	if err != nil {
		return err
	}

	if app.PermUrl == "" {
		if len(dbPerms) > 0 {
			// delete
			for _, dbPerm := range dbPerms {
				err = c.r.Perm.Delete(ctx, dbPerm.Id)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	respObj, err := c.fetchPerms(app.PermUrl)
	if err != nil {
		return err
	}

	var found bool

	for _, dbPerm := range dbPerms {
		found = false

		for _, freshPerm := range respObj.Perms {
			if freshPerm.Code == dbPerm.Code {
				found = true
				break
			}
		}

		if !found {
			// delete
			err = c.r.Perm.Delete(ctx, dbPerm.Id)
			if err != nil {
				return err
			}
		}
	}

	for _, freshPerm := range respObj.Perms {
		if freshPerm.Code == "" {
			continue
		}

		found = false

		for _, dbPerm := range dbPerms {
			if freshPerm.Code == dbPerm.Code {
				found = true

				if freshPerm.IsAll != dbPerm.IsAll || freshPerm.Dsc != dbPerm.Dsc {
					// update
					err = c.r.Perm.Update(ctx, dbPerm.Id, &entities.PermCUSt{
						IsAll: &freshPerm.IsAll,
						Dsc:   &freshPerm.Dsc,
					})
					if err != nil {
						return err
					}
				}

				break
			}
		}

		if !found {
			// create
			_, err = c.r.Perm.Create(ctx, &entities.PermCUSt{
				AppId:     &id,
				Code:      &freshPerm.Code,
				IsAll:     &freshPerm.IsAll,
				Dsc:       &freshPerm.Dsc,
				IsFetched: &cns.True,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *App) fetchPerms(permUrl string) (*entities.SystemGetPermsRepSt, error) {
	const fetchTimeout = 5 * time.Second

	var err error

	if permUrl == "" {
		return nil, errs.AppHasNotPermsUrl
	}

	result := &entities.SystemGetPermsRepSt{}

	httpClient := httpclient.New(c.r.lg, &httpc.OptionsSt{
		Client: &http.Client{
			Timeout:   fetchTimeout,
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		},
	})

	_, err = httpClient.Send(&httpc.OptionsSt{
		Method:   "GET",
		Uri:      permUrl,
		LogFlags: httpc.NoLogError,

		RepObj: result,
	})
	if err != nil {
		return nil, err
	}

	if result.Perms == nil {
		result.Perms = []*entities.SystemGetPermsItemSt{}
	}

	// filter
	filteredItems := make([]*entities.SystemGetPermsItemSt, 0, len(result.Perms))
	for _, item := range result.Perms {
		if item.Code != "" {
			filteredItems = append(filteredItems, item)
		}
	}
	result.Perms = filteredItems

	return result, nil
}
