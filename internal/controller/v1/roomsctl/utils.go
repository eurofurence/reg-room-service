package roomsctl

import (
	"context"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/google/uuid"
)

func validateRoomID(ctx context.Context, roomID string) error {
	if err := uuid.Validate(roomID); err != nil {
		return common.NewBadRequest(ctx, common.RoomIDInvalid, common.Details("you must specify a valid uuid"), err)
	}

	return nil
}
