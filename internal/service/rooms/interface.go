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
	FindMyRoom(ctx context.Context) (*modelsv1.Room, error)
}

type FindRoomParams struct {
	memberIDs []int64

	minSize uint
	maxSize uint

	minOccupants uint
	maxOccupants uint
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
