package dbrepo

import (
	"context"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/eurofurence/reg-room-service/internal/repository/database/historizeddb"
	"github.com/eurofurence/reg-room-service/internal/repository/database/mysqldb"
	"os"
	"time"

	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
)

var (
	activeRepository database.Repository

	inMemoryDBName  = "inmemdb"
	inMemoryAddress = "localhost"
	inMemoryPort    = "8806" // TODO read from config too

	dbServer *server.Server
)

func createInMemoryDB() (err error) {
	db := memory.NewDatabase(inMemoryDBName)
	db.EnablePrimaryKeyIndexes()

	engine := sqle.NewDefault(
		memory.NewDBProvider(
			db,
		))

	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%s", inMemoryAddress, inMemoryPort),
	}
	dbServer, err = server.NewDefaultServer(config, engine)
	return err
}

func startInMemoryDB() {
	if err := dbServer.Start(); err != nil {
		aulogging.Logger.NoCtx().Error().WithErr(err).Printf("failed to start inmemory database - BAILING OUT: %s", err.Error())
		os.Exit(1)
	}
}

func stopInMemoryDB() {
	_ = dbServer.Close()
}

func inMemoryConnectString() string {
	return fmt.Sprintf("no_user:@tcp(%s:%s)/%s", inMemoryAddress, inMemoryPort, inMemoryDBName)
}

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
		aulogging.Warn(ctx, "Starting inmemory database...")
		if err := createInMemoryDB(); err != nil {
			return err
		}

		go startInMemoryDB()
		time.Sleep(1 * time.Second)

		aulogging.Warn(ctx, "Opening inmemory database (not useful for production!)...")
		r = historizeddb.Create(mysqldb.Create(inMemoryConnectString()))
	}
	err := r.Open(ctx)
	SetRepository(r)
	return err
}

func Close(ctx context.Context) {
	aulogging.Info(ctx, "Closing database...")
	GetRepository().Close(ctx)
	SetRepository(nil)

	if dbServer != nil {
		stopInMemoryDB()
	}
}

func Migrate(ctx context.Context) error {
	aulogging.Info(ctx, "Migrating database...")
	return GetRepository().Migrate(ctx)
}
