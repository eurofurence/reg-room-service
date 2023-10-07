package rooms

import (
	"net/http"

	"github.com/eurofurence/reg-room-service/internal/controller"
	"github.com/eurofurence/reg-room-service/internal/web/common"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(router chi.Router, ctrl controller.Controller) {
	h := &Handler{
		ctrl: ctrl,
	}

	router.Route("/rooms", func(sr chi.Router) {
		initGetRoutes(sr, h)
		initPostRoutes(sr, h)
		initPutRoutes(sr, h)
		initDeleteRoutes(sr, h)
	})
}

func initGetRoutes(router chi.Router, h *Handler) {
	router.Method(
		http.MethodGet,
		"/",
		common.CreateHandler(
			h.ListRooms,
			h.ListRoomsRequest,
			h.ListRoomsResponse,
		),
	)

	router.Method(
		http.MethodGet,
		"/my",
		common.CreateHandler(
			h.FindMyRooom,
			h.FindMyRoomRequest,
			h.FindMyRoomResponse,
		),
	)
}

func initPostRoutes(router chi.Router, h *Handler) {

}

func initPutRoutes(router chi.Router, h *Handler) {

}

func initDeleteRoutes(router chi.Router, h *Handler) {

}
