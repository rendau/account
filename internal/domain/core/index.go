package core

import (
	"context"
	"sync"
	"time"

	"github.com/rendau/dop/adapters/cache"
	"github.com/rendau/dop/adapters/db"
	"github.com/rendau/dop/adapters/logger"
	"github.com/rendau/dop/adapters/sms"

	"github.com/rendau/account/internal/adapters/repo"
)

type JwtI interface {
	Create(ctx context.Context, sub string, ttl time.Duration, payload any) (string, error)
}

type St struct {
	lg         logger.Lite
	cache      cache.Cache
	db         db.RDBContextTransaction
	repo       repo.Repo
	jwts       JwtI
	sms        sms.Sms
	noSmsCheck bool
	testing    bool

	wg sync.WaitGroup

	Dic     *Dic
	System  *System
	Session *Session

	Config *Config
	Perm   *Perm
	Role   *Role
	Usr    *Usr
	App    *App
}

func New(
	lg logger.Lite,
	cache cache.Cache,
	db db.RDBContextTransaction,
	repo repo.Repo,
	jwts JwtI,
	sms sms.Sms,
	noSmsCheck bool,
	testing bool,
) *St {
	c := &St{
		lg:         lg,
		cache:      cache,
		db:         db,
		repo:       repo,
		jwts:       jwts,
		sms:        sms,
		noSmsCheck: noSmsCheck,
		testing:    testing,
	}

	c.Dic = NewDic(c)
	c.System = NewSystem(c)
	c.Session = NewSession(c)

	c.Config = NewConfig(c)
	c.Perm = NewPerm(c)
	c.Role = NewRole(c)
	c.Usr = NewUsr(c)
	c.App = NewApp(c)

	return c
}

func (c *St) WaitJobs() {
	c.wg.Wait()
}
