package groups

import (
	"context"
	"fmt"
	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"
	"github.com/eurofurence/reg-room-service/internal/web/common"
	"github.com/google/uuid"
	"net/http"
)

// Controller implements methods, which satisfy the endpoint format
// in the `common` package.
type Controller struct {
	ctrl groupservice.Service
}

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
