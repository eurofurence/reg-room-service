package dbrepo

import (
	"context"
	"os"
	"strings"

	aulogging "github.com/StephanHCB/go-autumn-logging"

	"github.com/eurofurence/reg-room-service/internal/repository/database"
	"github.com/eurofurence/reg-room-service/internal/repository/database/historizeddb"
	"github.com/eurofurence/reg-room-service/internal/repository/database/inmemorydb"
	"github.com/eurofurence/reg-room-service/internal/repository/database/mysqldb"
)

var activeRepository database.Repository

func SetRepository(repository database.Repository) {
	activeRepository = repository
}

func GetRepository() database.Repository {
	if activeRepository == nil {
		aulogging.Logger.NoCtx().Error().Print("You must Open() the database before using it. This is an error in your implementation.")
		os.Exit(1)
	}
	return activeRepository
}

func Open(ctx context.Context, variant string, mysqlConnectString string) error {
	var r database.Repository
	if variant == "mysql" {
		aulogging.Info(ctx, "Opening mysql database...")
		r = historizeddb.Create(mysqldb.Create(mysqlConnectString))
	} else {
		aulogging.Warn(ctx, "Opening inmemory database (not useful for production!)...")
		r = historizeddb.Create(inmemorydb.Create())
	}
	err := r.Open(ctx)
	SetRepository(r)
	return err
}

func Close(ctx context.Context) {
	aulogging.Info(ctx, "Closing database...")
	GetRepository().Close(ctx)
	SetRepository(nil)
}

func Migrate(ctx context.Context) error {
	aulogging.Info(ctx, "Migrating database...")
	return GetRepository().Migrate(ctx)
}

func MysqlConnectString(username string, password string, databaseName string, parameters []string) string {
	return username + ":" + password + "@" +
		databaseName + "?" + strings.Join(parameters, "&")
}
