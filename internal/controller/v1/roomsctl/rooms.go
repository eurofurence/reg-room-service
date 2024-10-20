package roomsctl

import (
	"github.com/eurofurence/reg-room-service/internal/application/web"
	roomservice "github.com/eurofurence/reg-room-service/internal/service/rooms"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Controller implements methods which satisfy the endpoint format
// in the `common` package.
type Controller struct {
	svc roomservice.Service
}

func InitRoutes(router chi.Router, svc roomservice.Service) {
	h := &Controller{
		svc: svc,
	}

	router.Route("/api/rest/v1/rooms", func(sr chi.Router) {
		initGetRoutes(sr, h)
		initPostRoutes(sr, h)
		initPutRoutes(sr, h)
		initDeleteRoutes(sr, h)
	})
}

func initGetRoutes(router chi.Router, h *Controller) {
	router.Method(
		http.MethodGet,
		"/",
		web.CreateHandler(
			h.ListRooms,
			h.ListRoomsRequest,
			h.ListRoomsResponse,
		),
	)

	router.Method(
		http.MethodGet,
		"/my",
		web.CreateHandler(
			h.FindMyRoom,
			h.FindMyRoomRequest,
			h.FindMyRoomResponse,
		),
	)

	router.Method(
		http.MethodGet,
		"/{uuid}",
		web.CreateHandler(
			h.GetRoomByID,
			h.GetRoomByIDRequest,
			h.GetRoomByIDResponse,
		),
	)
}

func initPostRoutes(router chi.Router, h *Controller) {
	router.Method(
		http.MethodPost,
		"/",
		web.CreateHandler(
			h.CreateRoom,
			h.CreateRoomRequest,
			h.CreateRoomResponse,
		),
	)

	router.Method(
		http.MethodPost,
		"/{uuid}/occupants/{badgenumber}",
		web.CreateHandler(
			h.AddToRoom,
			h.AddToRoomRequest,
			h.AddToRoomResponse,
		),
	)
}

func initPutRoutes(router chi.Router, h *Controller) {
	router.Method(
		http.MethodPut,
		"/{uuid}",
		web.CreateHandler(
			h.UpdateRoom,
			h.UpdateRoomRequest,
			h.UpdateRoomResponse,
		),
	)
}

func initDeleteRoutes(router chi.Router, h *Controller) {
	router.Method(
		http.MethodDelete,
		"/{uuid}",
		web.CreateHandler(
			h.DeleteRoom,
			h.DeleteRoomRequest,
			h.DeleteRoomResponse,
		),
	)

	router.Method(
		http.MethodDelete,
		"/{uuid}/occupants/{badgenumber}",
		web.CreateHandler(
			h.RemoveFromRoom,
			h.RemoveFromRoomRequest,
			h.RemoveFromRoomResponse,
		),
	)
}
