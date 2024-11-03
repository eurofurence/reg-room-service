package roomservice

import (
	"context"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/mailservice"
)

// Service defines the interface for the service function implementations for the room endpoints.
type Service interface {
	GetRoomByID(ctx context.Context, roomID string) (*modelsv1.Room, error)
	CreateRoom(ctx context.Context, room *modelsv1.RoomCreate) (string, error)
	UpdateRoom(ctx context.Context, room *modelsv1.Room) error
	DeleteRoom(ctx context.Context, roomID string) error

	AddOccupantToRoom(ctx context.Context, roomID string, badgeNumber int64) error
	RemoveOccupantFromRoom(ctx context.Context, roomID string, badgeNumber int64) error

	FindRooms(ctx context.Context, params *FindRoomParams) ([]*modelsv1.Room, error)
	// FindMyRoom looks up the room the currently logged-in user is in.
	//
	// This works for admins just like for normal users, returning their room,
	// but will fail for requests using an API Token (no currently logged-in user available).
	//
	// Only finds rooms that have the "final" flag, and only works if the user's registration has
	// attending status.
	FindMyRoom(ctx context.Context) (*modelsv1.Room, error)
}

type FindRoomParams struct {
	MemberIDs []int64 // empty list or nil means no condition

	MinSize uint // 0 means no condition
	MaxSize uint // 0 means no condition

	MinOccupants uint // 0 means no condition
	MaxOccupants int  // -1 means no condition, 0 means search for empty rooms only
}

func New(db database.Repository, attsrv attendeeservice.AttendeeService, mailsrv mailservice.MailService) Service {
	return &roomService{
		DB:      db,
		AttSrv:  attsrv,
		MailSrv: mailsrv,
	}
}

type roomService struct {
	DB      database.Repository
	AttSrv  attendeeservice.AttendeeService
	MailSrv mailservice.MailService
}
