package groups

import (
	"context"
	"fmt"
	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/web/common"
	"github.com/google/uuid"
	"net/http"
)

func validateGroupID(ctx context.Context, w http.ResponseWriter, groupID string) error {
	if err := uuid.Validate(groupID); err != nil {
		common.SendHTTPStatusErrorResponse(
			ctx,
			w,
			apierrors.NewBadRequest(common.GroupIDInvalidMessage, fmt.Sprintf("%q is not a vailid UUID", groupID)))

		return err
	}

	return nil
}
