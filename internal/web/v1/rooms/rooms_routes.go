package rooms

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/eurofurence/reg-room-service/internal/controller"
	"github.com/eurofurence/reg-room-service/internal/web/common"
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

	router.Method(
		http.MethodGet,
		"/{uuid}",
		common.CreateHandler(
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
		common.CreateHandler(
			h.CreateRoom,
			h.CreateRoomRequest,
			h.CreateRoomResponse,
		),
	)

	router.Method(
		http.MethodPost,
		"/{uuid}/individuals/{badgenumber}",
		common.CreateHandler(
			h.AddRoomMember,
			h.AddRoomMemberRequest,
			h.AddRoomMemberResponse,
		),
	)

	router.Method(
		http.MethodPost,
		"/{uuid}/groups/{groupid}",
		common.CreateHandler(
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
		common.CreateHandler(
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
		common.CreateHandler(
			h.DeleteRoom,
			h.DeleteRoomRequest,
			h.DeleteRoomResponse,
		),
	)

	router.Method(
		http.MethodDelete,
		"/{uuid}/individuals/{badgenumber}",
		common.CreateHandler(
			h.DeleteRoomMember,
			h.DeleteRoomMemberRequest,
			h.DeleteRoomMemberResponse,
		),
	)

	router.Method(
		http.MethodPost,
		"/{uuid}/groups/{groupid}",
		common.CreateHandler(
			h.DeleteGroup,
			h.DeleteGroupRequest,
			h.DeleteGroupResponse,
		),
	)
}
