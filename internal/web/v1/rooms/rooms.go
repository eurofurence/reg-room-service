package rooms

import "github.com/eurofurence/reg-room-service/internal/controller"

// Handler implements methods, which satisfy the endpoint format
type Handler struct {
	ctrl controller.Controller
}
