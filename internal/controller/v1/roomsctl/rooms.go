package roomsctl

import (
	"github.com/eurofurence/reg-room-service/internal/application/web"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Handler implements methods, which satisfy the endpoint format.
type Handler struct{}

func InitRoutes(router chi.Router) {
	h := &Handler{}

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
			h.FindMyRooom,
			h.FindMyRoomRequest,
			h.FindMyRoomResponse,
		),
	)

	router.Method(
		http.MethodGet,
		"/{uuid}",
		web.CreateHandler(
			h.FindRoomByUUID,
			h.FindRoomByUUIDRequest,
			h.FindRoomByUUIDResponse,
		),
	)
}

func initPostRoutes(router chi.Router, h *Handler) {
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
		"/{uuid}/individuals/{badgenumber}",
		web.CreateHandler(
			h.AddRoomMember,
			h.AddRoomMemberRequest,
			h.AddRoomMemberResponse,
		),
	)

	router.Method(
		http.MethodPost,
		"/{uuid}/groups/{groupid}",
		web.CreateHandler(
			h.AddGroup,
			h.AddGroupRequest,
			h.AddGroupResponse,
		),
	)
}

func initPutRoutes(router chi.Router, h *Handler) {
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

func initDeleteRoutes(router chi.Router, h *Handler) {
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
		"/{uuid}/individuals/{badgenumber}",
		web.CreateHandler(
			h.DeleteRoomMember,
			h.DeleteRoomMemberRequest,
			h.DeleteRoomMemberResponse,
		),
	)

	router.Method(
		http.MethodDelete,
		"/{uuid}/groups/{groupid}",
		web.CreateHandler(
			h.DeleteGroup,
			h.DeleteGroupRequest,
			h.DeleteGroupResponse,
		),
	)
}
