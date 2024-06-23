package app

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	auzerolog "github.com/StephanHCB/go-autumn-logging-zerolog"
	"github.com/eurofurence/reg-room-service/internal/config"
	"github.com/eurofurence/reg-room-service/internal/repository/database/dbrepo"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"github.com/eurofurence/reg-room-service/internal/web/common"
	"github.com/eurofurence/reg-room-service/internal/web/server"
	"github.com/rs/zerolog"
)

type Params struct {
	configFilePath string
	migrateDB      bool
}

func NewParams(configFile string, migrateDB bool) Params {
	return Params{
		configFilePath: configFile,
		migrateDB:      migrateDB,
	}
}

type Application struct {
	Params Params
}

func New(params Params) *Application {
	return &Application{
		Params: params,
	}
}

func (a *Application) Run() error {
	conf, err := config.UnmarshalFromYamlConfiguration(a.Params.configFilePath)
	setupLogging(conf)
	ctx := auzerolog.AddLoggerToCtx(context.Background())
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "failed to load configuration - bailing out: %s", err.Error())
		return err
	}

	if conf.Database.Use == config.Mysql {
		connectString := dbrepo.MysqlConnectString(conf.Database.Username, conf.Database.Password, conf.Database.Database, conf.Database.Parameters)
		err := dbrepo.Open(ctx, string(config.Mysql), connectString)
		if err != nil {
			aulogging.ErrorErrf(ctx, err, "failed to set up database connection - bailing out: %s", err.Error())
			return err
		}
	} else {
		err := dbrepo.Open(ctx, string(config.Inmemory), "")
		if err != nil {
			aulogging.ErrorErrf(ctx, err, "failed to set up inmemory database - bailing out: %s", err.Error())
			return err
		}
	}

	attsrv, err := attendeeservice.New(conf.Service.AttendeeServiceURL)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "failed to set up attendee service client - bailing out: %s", err.Error())
		return err
	}

	srv := server.NewServer(conf, context.Background())
	err = srv.Serve(dbrepo.GetRepository(), attsrv)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "failure during serve phase - shutting down: %s", err.Error())
		return err
	}

	aulogging.Info(ctx, "done serving web requests")
	return nil
}

const applicationName = "room-service"

func setupLogging(confOrNil *config.Config) {
	useEcsLogging := confOrNil != nil && confOrNil.Logging.Style == config.ECS
	severity := "INFO"
	if confOrNil != nil && confOrNil.Logging.Severity != "" {
		severity = confOrNil.Logging.Severity
	}

	aulogging.RequestIdRetriever = common.GetRequestID
	if useEcsLogging {
		auzerolog.SetupJsonLogging(applicationName)
		zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z"
	} else {
		aulogging.DefaultRequestIdValue = "00000000"
		auzerolog.SetupPlaintextLogging()
	}

	switch severity {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
