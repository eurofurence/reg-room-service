package mysqldb

import (
	"context"
	"errors"
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

func (r *MysqlRepository) Migrate(_ context.Context) error {
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

const groupDesc = "group"

func (r *MysqlRepository) AddGroup(ctx context.Context, g *entity.Group) (string, error) {
	err := add[entity.Group](ctx, r.db, g, groupDesc)
	return g.ID, err
}

func (r *MysqlRepository) UpdateGroup(ctx context.Context, g *entity.Group) error {
	return update[entity.Group](ctx, r.db, g, groupDesc)
}

func (r *MysqlRepository) GetGroupByID(ctx context.Context, id string) (*entity.Group, error) {
	return getByID[entity.Group](ctx, r.db, id, groupDesc)
}

func (r *MysqlRepository) SoftDeleteGroupByID(ctx context.Context, id string) error {
	return softDeleteByID[entity.Group](ctx, r.db, id, groupDesc)
}

func (r *MysqlRepository) UndeleteGroupByID(ctx context.Context, id string) error {
	return undeleteByID[entity.Group](ctx, r.db, id, groupDesc)
}

func (r *MysqlRepository) NewEmptyGroupMembership(_ context.Context, groupID string, attendeeID uint) *entity.GroupMember {
	var m entity.GroupMember
	m.ID = attendeeID
	m.GroupID = groupID
	m.IsInvite = true // default to invite because that's the usual starting point
	return &m
}

const groupMembershipDesc = "group membership"

func (r *MysqlRepository) GetGroupMembershipByAttendeeID(ctx context.Context, attendeeID uint) (*entity.GroupMember, error) {
	var m entity.GroupMember
	m.ID = attendeeID
	return getMembershipByAttendeeID[entity.GroupMember](ctx, r.db, attendeeID, &m, groupMembershipDesc)
}

func (r *MysqlRepository) GetGroupMembersByGroupID(ctx context.Context, groupID string) ([]entity.GroupMember, error) {
	return selectMembersBy[entity.GroupMember](ctx, r.db, &entity.GroupMember{GroupID: groupID}, groupMembershipDesc)
}

func (r *MysqlRepository) AddGroupMembership(ctx context.Context, gm *entity.GroupMember) error {
	return addMembership[entity.GroupMember](ctx, r.db, gm, groupMembershipDesc)
}

func (r *MysqlRepository) UpdateGroupMembership(ctx context.Context, gm *entity.GroupMember) error {
	return updateMembership[entity.GroupMember](ctx, r.db, gm, groupMembershipDesc)
}

func (r *MysqlRepository) DeleteGroupMembership(ctx context.Context, attendeeID uint) error {
	return deleteMembership[entity.GroupMember](ctx, r.db, attendeeID, groupMembershipDesc)
}

const roomDesc = "room"

func (r *MysqlRepository) AddRoom(ctx context.Context, room *entity.Room) (string, error) {
	err := add[entity.Room](ctx, r.db, room, roomDesc)
	return room.ID, err
}

func (r *MysqlRepository) UpdateRoom(ctx context.Context, room *entity.Room) error {
	return update[entity.Room](ctx, r.db, room, roomDesc)
}

func (r *MysqlRepository) GetRoomByID(ctx context.Context, id string) (*entity.Room, error) {
	return getByID[entity.Room](ctx, r.db, id, roomDesc)
}

func (r *MysqlRepository) SoftDeleteRoomByID(ctx context.Context, id string) error {
	return softDeleteByID[entity.Room](ctx, r.db, id, roomDesc)
}

func (r *MysqlRepository) UndeleteRoomByID(ctx context.Context, id string) error {
	return undeleteByID[entity.Room](ctx, r.db, id, roomDesc)
}

const roomMembershipDesc = "room membership"

func (r *MysqlRepository) NewEmptyRoomMembership(_ context.Context, roomID string, attendeeID uint) *entity.RoomMember {
	var m entity.RoomMember
	m.ID = attendeeID
	m.RoomID = roomID
	return &m
}

func (r *MysqlRepository) GetRoomMembershipByAttendeeID(ctx context.Context, attendeeID uint) (*entity.RoomMember, error) {
	var m entity.RoomMember
	m.ID = attendeeID
	return getMembershipByAttendeeID[entity.RoomMember](ctx, r.db, attendeeID, &m, roomMembershipDesc)
}

func (r *MysqlRepository) GetRoomMembersByRoomID(ctx context.Context, roomID string) ([]entity.RoomMember, error) {
	return selectMembersBy[entity.RoomMember](ctx, r.db, &entity.RoomMember{RoomID: roomID}, roomMembershipDesc)
}

func (r *MysqlRepository) AddRoomMembership(ctx context.Context, rm *entity.RoomMember) error {
	return addMembership[entity.RoomMember](ctx, r.db, rm, roomMembershipDesc)
}

func (r *MysqlRepository) UpdateRoomMembership(ctx context.Context, rm *entity.RoomMember) error {
	return updateMembership[entity.RoomMember](ctx, r.db, rm, roomMembershipDesc)
}

func (r *MysqlRepository) DeleteRoomMembership(ctx context.Context, attendeeID uint) error {
	return deleteMembership[entity.RoomMember](ctx, r.db, attendeeID, roomMembershipDesc)
}

func (r *MysqlRepository) RecordHistory(ctx context.Context, h *entity.History) error {
	err := r.db.Create(h).Error
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("mysql error during history entry insert: %s", err.Error())
	}
	return err
}

// generics to reduce repetitions

type anyMemberCollection interface {
	entity.Group | entity.Room
}

func add[E anyMemberCollection](
	ctx context.Context,
	db *gorm.DB,
	c *E,
	logDescription string,
) error {
	err := db.Create(c).Error
	if err != nil {
		aulogging.WarnErrf(ctx, err, "mysql error during %s insert: %s", logDescription, err.Error())
	}
	return err
}

func update[E anyMemberCollection](
	ctx context.Context,
	db *gorm.DB,
	c *E,
	logDescription string,
) error {
	// does not allow updating deleted groups/rooms, use .Unscoped to allow
	err := db.Save(c).Error
	if err != nil {
		aulogging.WarnErrf(ctx, err, "mysql error during %s update: %s", logDescription, err.Error())
	}
	return err
}

func getByID[E anyMemberCollection](
	ctx context.Context,
	db *gorm.DB,
	id string,
	logDescription string,
) (*E, error) {
	var g E
	// allow reading deleted so history and undelete work
	err := db.Unscoped().First(&g, id).Error
	if err != nil {
		aulogging.InfoErrf(ctx, err, "mysql error during %s select - might be ok: %s", logDescription, err.Error())
	}
	return &g, err
}

func softDeleteByID[E anyMemberCollection](
	ctx context.Context,
	db *gorm.DB,
	id string,
	logDescription string,
) error {
	var g E
	err := db.First(&g, id).Error
	if err != nil {
		aulogging.WarnErrf(ctx, err, "mysql error during %s soft delete - %s not found: %s", logDescription, logDescription, err.Error())
		return err
	}
	err = db.Delete(&g).Error
	if err != nil {
		aulogging.WarnErrf(ctx, err, "mysql error during %s soft delete - deletion failed: %s", logDescription, err.Error())
		return err
	}
	return nil
}

func undeleteByID[E anyMemberCollection](
	ctx context.Context,
	db *gorm.DB,
	id string,
	logDescription string,
) error {
	var g E
	err := db.Unscoped().First(&g, id).Error
	if err != nil {
		aulogging.WarnErrf(ctx, err, "mysql error during %s undelete - %s not found: %s", logDescription, logDescription, err.Error())
		return err
	}
	err = db.Unscoped().Model(&g).Where("id", id).Update("deleted_at", nil).Error
	if err != nil {
		aulogging.WarnErrf(ctx, err, "mysql error during %s undelete: %s", logDescription, err.Error())
		return err
	}
	return nil
}

type anyMembership interface {
	entity.GroupMember | entity.RoomMember
}

func getMembershipByAttendeeID[E anyMembership](
	ctx context.Context,
	db *gorm.DB,
	attendeeID uint,
	defaultValue *E,
	logDescription string,
) (*E, error) {
	var m E
	err := db.First(&m, attendeeID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			aulogging.Infof(ctx, "no %s for attendee id %d - might be ok", logDescription, attendeeID)
			return defaultValue, err
		} else {
			aulogging.WarnErrf(ctx, err, "mysql error during %s select - not record not found: %s", logDescription, err.Error())
			return defaultValue, err
		}
	}
	return &m, err
}

func selectMembersBy[E anyMembership](
	ctx context.Context,
	db *gorm.DB,
	condition *E,
	logDescription string,
) ([]E, error) {
	var table E
	rows, err := db.Model(&table).Where(condition).Rows()
	if err != nil {
		aulogging.WarnErrf(ctx, err, "mysql error during %s select: %s", logDescription, err.Error())
		return make([]E, 0), err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			aulogging.WarnErrf(ctx, err, "mysql error during %s result set close: %s", logDescription, err.Error())
		}
	}()

	result := make([]E, 0)
	for rows.Next() {
		var sc E
		err := db.ScanRows(rows, &sc)
		if err != nil {
			aulogging.WarnErrf(ctx, err, "mysql error during %s read: %s", logDescription, err.Error())
			return make([]E, 0), err
		}

		result = append(result, sc)
	}

	return result, nil
}

func addMembership[E anyMembership](
	ctx context.Context,
	db *gorm.DB,
	m *E,
	logDescription string,
) error {
	err := db.Create(m).Error
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("mysql error during %s insert: %s", logDescription, err.Error())
	}
	return err
}

func updateMembership[E anyMembership](
	ctx context.Context,
	db *gorm.DB,
	m *E,
	logDescription string,
) error {
	err := db.Save(m).Error
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("mysql error during %s update: %s", logDescription, err.Error())
	}
	return err
}

func deleteMembership[E anyMembership](
	ctx context.Context,
	db *gorm.DB,
	id uint,
	logDescription string,
) error {
	var m E
	err := db.First(&m, id).Error
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("mysql error during %s delete - not found: %s", logDescription, err.Error())
		return err
	}
	err = db.Delete(&m).Error
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("mysql error during %s delete - deletion failed: %s", logDescription, err.Error())
		return err
	}
	return nil
}
