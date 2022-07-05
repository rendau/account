package tests

import (
	"log"
	"os"
	"testing"

	"github.com/rendau/account/internal/adapters/repo/pg"
	"github.com/rendau/account/internal/domain/core"
	"github.com/rendau/account/internal/domain/usecases"
	"github.com/rendau/dop/adapters/cache/mem"
	dopDbPg "github.com/rendau/dop/adapters/db/pg"
	dopJwtMock "github.com/rendau/dop/adapters/jwt/mock"
	dopLoggerZap "github.com/rendau/dop/adapters/logger/zap"
	dopSmsMock "github.com/rendau/dop/adapters/sms/mock"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	var err error

	viper.SetConfigFile("conf.yml")
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	app.lg = dopLoggerZap.New(
		"info",
		true,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer app.lg.Sync()

	app.cache = mem.New()

	app.db, err = dopDbPg.New(true, app.lg, dopDbPg.OptionsSt{
		Dsn: viper.GetString("pg_dsn"),
	})
	if err != nil {
		app.lg.Fatal(err)
	}

	app.repo = pg.New(app.db, app.lg)

	app.ucs = usecases.New(
		app.lg,
		app.db,
	)

	app.jwts = dopJwtMock.New(app.lg, true)

	app.sms = dopSmsMock.New(app.lg, true)

	app.core = core.New(
		app.lg,
		app.cache,
		app.db,
		app.repo,
		app.jwts,
		app.sms,
		false,
		true,
	)

	app.ucs.SetCore(app.core)

	resetDb()

	// Start tests
	code := m.Run()

	os.Exit(code)
}
