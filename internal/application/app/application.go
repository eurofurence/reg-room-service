package app

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	auzerolog "github.com/StephanHCB/go-autumn-logging-zerolog"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/application/server"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/repository/database/dbrepo"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/authservice"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/mailservice"
	groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"
	roomservice "github.com/eurofurence/reg-room-service/internal/service/rooms"
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
	// config and logging

	conf, err := config.UnmarshalFromYamlConfiguration(a.Params.configFilePath)
	setupLogging(conf)
	ctx := auzerolog.AddLoggerToCtx(context.Background())
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "failed to load configuration - bailing out: %s", err.Error())
		return err
	}
	aulogging.Info(ctx, "configuration file successfully loaded")

	aulogging.Info(ctx, "adding configuration defaults")
	conf.AddDefaults()
	aulogging.Info(ctx, "applying environment variable overrides")
	conf.ApplyEnvironmentOverrides()
	aulogging.Info(ctx, "validating configuration")
	err = conf.Validate()
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "failed to validate configuration - bailing out: %s", err.Error())
		return err
	}

	// repos

	connectString := dbrepo.MysqlConnectString(conf.Database.Username, conf.Database.Password, conf.Database.Database, conf.Database.Parameters)
	if err := dbrepo.Open(ctx, string(conf.Database.Use), connectString); err != nil {
		aulogging.ErrorErrf(ctx, err, "failed to set up database connection - bailing out: %s", err.Error())
		return err
	}

	if a.Params.migrateDB {
		err := dbrepo.Migrate(ctx)
		if err != nil {
			aulogging.ErrorErrf(ctx, err, "failed to migrate database - bailing out: %s", err.Error())
			return err
		}
	}

	attRepo, err := attendeeservice.New(conf.Service.AttendeeServiceURL)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "failed to set up attendee service client - bailing out: %s", err.Error())
		return err
	}

	_, err = authservice.New(conf.Service.AuthServiceURL, conf.Security)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "failed to set up auth service client - bailing out: %s", err.Error())
		return err
	}

	mailRepo, err := mailservice.New(conf.Service.MailServiceURL, conf.Security.Fixed.API)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "failed to set up mail service client - bailing out: %s", err.Error())
		return err
	}

	dbRepo := dbrepo.GetRepository()

	// services

	groupSvc := groupservice.New(dbRepo, attRepo, mailRepo)
	roomSvc := roomservice.New(dbRepo, attRepo, mailRepo)

	// controllers wired in server because no instances, just routes

	srv := server.New(conf, context.Background(), groupSvc, roomSvc)
	err = srv.Serve()
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
