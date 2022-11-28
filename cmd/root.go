package cmd

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/rendau/account/docs"
	"github.com/rendau/account/internal/adapters/repo"
	"github.com/rendau/account/internal/adapters/repo/pg"
	"github.com/rendau/account/internal/adapters/server/rest"
	"github.com/rendau/account/internal/domain/core"
	"github.com/rendau/account/internal/domain/usecases"
	dopCache "github.com/rendau/dop/adapters/cache"
	dopCacheMem "github.com/rendau/dop/adapters/cache/mem"
	dopCacheRedis "github.com/rendau/dop/adapters/cache/redis"
	"github.com/rendau/dop/adapters/client/httpc"
	"github.com/rendau/dop/adapters/client/httpc/httpclient"
	dopDbPg "github.com/rendau/dop/adapters/db/pg"
	dopJwt "github.com/rendau/dop/adapters/jwt"
	dopJwts "github.com/rendau/dop/adapters/jwt/jwts"
	dopJwtMock "github.com/rendau/dop/adapters/jwt/mock"
	dopLoggerZap "github.com/rendau/dop/adapters/logger/zap"
	dopServerHttps "github.com/rendau/dop/adapters/server/https"
	dopSms "github.com/rendau/dop/adapters/sms"
	dopSmsMock "github.com/rendau/dop/adapters/sms/mock"
	dopSmss "github.com/rendau/dop/adapters/sms/smss"
	"github.com/rendau/dop/dopTools"
)

func Execute() {
	var err error

	app := struct {
		lg         *dopLoggerZap.St
		cache      dopCache.Cache
		db         *dopDbPg.St
		repo       repo.Repo
		jwts       dopJwt.Jwt
		sms        dopSms.Sms
		core       *core.St
		ucs        *usecases.St
		restApiSrv *dopServerHttps.St
	}{}

	confLoad()

	app.lg = dopLoggerZap.New(conf.LogLevel, conf.Debug)

	if conf.RedisUrl == "" {
		app.cache = dopCacheMem.New()
	} else {
		app.cache = dopCacheRedis.New(
			app.lg,
			conf.RedisUrl,
			conf.RedisPsw,
			conf.RedisDb,
			conf.RedisKeyPrefix,
		)
	}

	app.db, err = dopDbPg.New(conf.Debug, app.lg, dopDbPg.OptionsSt{
		Dsn: conf.PgDsn,
	})
	if err != nil {
		app.lg.Fatal(err)
	}

	app.repo = pg.New(app.db, app.lg)

	app.ucs = usecases.New(app.lg, app.db)

	if conf.MsJwtsUrl == "" {
		app.jwts = dopJwtMock.New(app.lg, false)
	} else {
		app.jwts = dopJwts.New(
			httpclient.New(app.lg, &httpc.OptionsSt{
				Client: &http.Client{
					Timeout: 10 * time.Second,
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					},
				},
				Uri:       conf.MsJwtsUrl,
				LogPrefix: "JWT: ",
			}),
		)
	}

	if conf.MsSmsUrl == "" {
		app.sms = dopSmsMock.New(app.lg, false)
	} else {
		app.sms = dopSmss.New(
			httpclient.New(app.lg, &httpc.OptionsSt{
				Client: &http.Client{
					Timeout: 10 * time.Second,
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					},
				},
				Uri:       conf.MsSmsUrl,
				LogPrefix: "SMS: ",
			}),
		)
	}

	app.core = core.New(
		app.lg,
		app.cache,
		app.db,
		app.repo,
		app.jwts,
		app.sms,
		conf.NoSmsCheck,
		false,
	)

	app.ucs.SetCore(app.core)

	docs.SwaggerInfo.Host = conf.SwagHost
	docs.SwaggerInfo.BasePath = conf.SwagBasePath
	docs.SwaggerInfo.Schemes = []string{conf.SwagSchema}
	docs.SwaggerInfo.Title = "Account service"

	// START

	app.lg.Infow("Starting")

	app.restApiSrv = dopServerHttps.Start(
		conf.HttpListen,
		rest.GetHandler(
			app.lg,
			app.ucs,
			conf.HttpCors,
		),
		app.lg,
	)

	var exitCode int

	select {
	case <-dopTools.StopSignal():
	case <-app.restApiSrv.Wait():
		exitCode = 1
	}

	// STOP

	app.lg.Infow("Shutting down...")

	if !app.restApiSrv.Shutdown(20 * time.Second) {
		exitCode = 1
	}

	app.lg.Infow("Wait routines...")

	app.core.WaitJobs()

	app.lg.Infow("Exit")

	os.Exit(exitCode)
}
