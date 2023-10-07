package mysqldb

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/entity"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type MysqlRepository struct {
	db            *gorm.DB
	connectString string
	Now           func() time.Time
}

func Create(connectString string) database.Repository {
	return &MysqlRepository{
		Now:           time.Now,
		connectString: connectString,
	}
}

func (r *MysqlRepository) Open(ctx context.Context) error {
	gormConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "room_",
		},
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(mysql.Open(r.connectString), &gormConfig)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "failed to open mysql connection: %s", err.Error())
		return err
	}

	sqlDb, err := db.DB()
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "failed to configure mysql connection: %s", err.Error())
		return err
	}

	// see https://making.pusher.com/production-ready-connection-pooling-in-go/
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetMaxIdleConns(50)
	sqlDb.SetConnMaxLifetime(time.Minute * 10)

	r.db = db
	return nil
}

func (r *MysqlRepository) Close(_ context.Context) {
	// no more db close in gorm v2
}

func (r *MysqlRepository) Migrate(ctx context.Context) error {
	err := r.db.AutoMigrate(
		&entity.Group{},
		&entity.GroupMember{},
		&entity.Room{},
		&entity.RoomMember{},
	)
	if err != nil {
		aulogging.Logger.NoCtx().Error().WithErr(err).Printf("failed to migrate mysql db: %s", err.Error())
		return err
	}
	return nil
}
