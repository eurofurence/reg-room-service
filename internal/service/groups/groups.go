package groupservice

import (
	"context"
	"errors"
	"fmt"
	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/entity"
	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
	"gorm.io/gorm"
	"strings"
)

type Service interface {
	GetGroupByID(ctx context.Context, groupID string) (*modelsv1.Group, error)
	CreateGroup(ctx context.Context, group modelsv1.GroupCreate) (string, error)
}

func NewService(repository database.Repository) Service {
	return &groupService{DB: repository}
}

type groupService struct {
	DB database.Repository
}

func deref[T any](ptr *T) T {
	var def T
	if ptr == nil {
		return def
	}

	return *ptr
}

func ptr[T any](val T) *T {
	return &val
}

func (g *groupService) GetGroupByID(ctx context.Context, groupID string) (*modelsv1.Group, error) {
	grp, err := g.DB.GetGroupByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierrors.NewNotFound("unable to find record", fmt.Sprintf("no record found for id %q", groupID))
		}

		return nil, apierrors.NewInternalServerError("something went wrong", err.Error())
	}

	return &modelsv1.Group{
		ID:          grp.ID,
		Name:        grp.Name,
		Flags:       strings.Split(grp.Flags, ","),
		Comments:    &grp.Comments,
		MaximumSize: ptr(int32(grp.MaximumSize)),
		Owner:       int32(grp.Owner),
		Members: []modelsv1.Member{
			{
				ID:       42,
				Nickname: "hardcoded", // TODO
			},
		},
		Invites: nil,
	}, nil
}

func (g *groupService) CreateGroup(ctx context.Context, group modelsv1.GroupCreate) (string, error) {
	return g.DB.AddGroup(ctx, &entity.Group{
		Name:        group.Name,
		Flags:       fmt.Sprintf(",%s,", strings.Join(group.Flags, ",")),
		Comments:    deref(group.Comments),
		MaximumSize: 6,  // TODO add from config
		Owner:       42, // TODO read from attendee service (or passed in by admin)
	})
}
