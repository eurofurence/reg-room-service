package roomsctl

import (
	"context"
	"fmt"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/google/uuid"
	"net/url"
)

func validateRoomID(ctx context.Context, roomID string) error {
	if err := uuid.Validate(roomID); err != nil {
		return common.NewBadRequest(ctx, common.RoomIDInvalid, common.Details(fmt.Sprintf("'%s' is not a valid UUID", url.PathEscape(roomID))), err)
	}

	return nil
}
