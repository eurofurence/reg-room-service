package groupsctl

import (
	"context"
	"fmt"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/google/uuid"
	"net/url"
)

func validateGroupID(ctx context.Context, groupID string) error {
	if err := uuid.Validate(groupID); err != nil {
		return common.NewBadRequest(ctx, common.GroupIDInvalid, common.Details(fmt.Sprintf("'%s' is not a valid UUID", url.PathEscape(groupID))), err)
	}

	return nil
}
