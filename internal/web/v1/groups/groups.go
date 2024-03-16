package groups

import groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"

// Controller implements methods, which satisfy the endpoint format
// in the `common` package.
type Controller struct {
	ctrl groupservice.Service
}
