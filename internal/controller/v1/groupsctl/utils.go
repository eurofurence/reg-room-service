package groupsctl

import (
	"context"
	"fmt"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/application/web"
	"github.com/google/uuid"
	"net/http"
	"net/url"
)

func validateGroupID(ctx context.Context, w http.ResponseWriter, groupID string) error {
	if err := uuid.Validate(groupID); err != nil {
		web.SendErrorResponse(ctx, w,
			common.NewBadRequest(ctx, common.GroupIDInvalid, common.Details(fmt.Sprintf("'%s' is not a valid UUID", url.PathEscape(groupID)))),
		)

		return err
	}

	return nil
}
