package groups

import (
	"github.com/eurofurence/reg-room-service/internal/controller"
)

// Handler implements methods, which satisfy the endpoint format
// in the `common` package.
type Handler struct {
	ctrl controller.Controller
}
