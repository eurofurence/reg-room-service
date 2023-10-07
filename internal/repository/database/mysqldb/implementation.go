package mysqldb

import (
	"context"
	"time"

	aulogging "github.com/StephanHCB/go-autumn-logging"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/eurofurence/reg-room-service/internal/entity"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
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
		&entity.History{},
		&entity.Room{},
		&entity.RoomMember{},
	)
	if err != nil {
		aulogging.Logger.NoCtx().Error().WithErr(err).Printf("failed to migrate mysql db: %s", err.Error())
		return err
	}
	return nil
}

func (r *MysqlRepository) AddGroup(ctx context.Context, g *entity.Group) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) UpdateGroup(ctx context.Context, g *entity.Group) error {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) GetGroupByID(ctx context.Context, id string) (*entity.Group, error) {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) SoftDeleteGroupByID(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) UndeleteGroupByID(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) NewEmptyGroupMembership(ctx context.Context, groupID string, attendeeID uint) *entity.GroupMember {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) GetGroupMembershipByAttendeeID(ctx context.Context, attendeeID uint) (*entity.GroupMember, error) {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) GetGroupMembersByGroupID(ctx context.Context, groupID string) ([]*entity.GroupMember, error) {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) AddGroupMembership(ctx context.Context, gm *entity.GroupMember) error {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) UpdateGroupMembership(ctx context.Context, gm *entity.GroupMember) error {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) DeleteGroupMembership(ctx context.Context, attendeeID uint) error {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) AddRoom(ctx context.Context, g *entity.Room) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) UpdateRoom(ctx context.Context, g *entity.Room) error {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) GetRoomByID(ctx context.Context, id string) (*entity.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) SoftDeleteRoomByID(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) UndeleteRoomByID(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) NewEmptyRoomMembership(ctx context.Context, roomID string, attendeeID uint) *entity.RoomMember {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) GetRoomMembershipByAttendeeID(ctx context.Context, attendeeID uint) (*entity.RoomMember, error) {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) GetRoomMembersByRoomID(ctx context.Context, roomID string) ([]*entity.RoomMember, error) {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) AddRoomMembership(ctx context.Context, gm *entity.RoomMember) error {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) UpdateRoomMembership(ctx context.Context, gm *entity.RoomMember) error {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) DeleteRoomMembership(ctx context.Context, attendeeID uint) error {
	//TODO implement me
	panic("implement me")
}

func (r *MysqlRepository) RecordHistory(ctx context.Context, h *entity.History) error {
	//TODO implement me
	panic("implement me")
}
