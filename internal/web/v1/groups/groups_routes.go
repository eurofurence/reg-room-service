package groups

import (
	"net/http"

	groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"

	"github.com/go-chi/chi/v5"

	"github.com/eurofurence/reg-room-service/internal/web/common"
)

func InitRoutes(router chi.Router, ctrl groupservice.Service) {
	h := &Controller{
		ctrl: ctrl,
	}

	router.Route("/api/rest/v1/groups", func(sr chi.Router) {
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
		common.CreateHandler(
			h.ListGroups,
			h.ListGroupsRequest,
			h.ListGroupsResponse,
		),
	)

	router.Method(
		http.MethodGet,
		"/my",
		common.CreateHandler(
			h.FindMyGroup,
			h.FindMyGroupRequest,
			h.FindMyGroupResponse,
		),
	)

	router.Method(
		http.MethodGet,
		"/{uuid}",
		common.CreateHandler(
			h.FindGroupByID,
			h.FindGroupByIDRequest,
			h.FindGroupByIDResponse,
		),
	)
}

func initPostRoutes(router chi.Router, h *Controller) {
	router.Method(
		http.MethodPost,
		"/",
		common.CreateHandler(
			h.CreateGroup,
			h.CreateGroupRequest,
			h.CreateGroupResponse,
		),
	)

	router.Method(
		http.MethodPost,
		"/{uuid}/members/{badgenumber}",
		common.CreateHandler(
			h.AddMemberToGroup,
			h.AddMemberToGroupRequest,
			h.AddMemberToGroupResponse,
		),
	)
}

func initPutRoutes(router chi.Router, h *Controller) {
	router.Method(
		http.MethodPut,
		"/{uuid}",
		common.CreateHandler(
			h.UpdateGroup,
			h.UpdateGroupRequest,
			h.UpdateGroupResponse,
		),
	)
}

func initDeleteRoutes(router chi.Router, h *Controller) {
	router.Method(
		http.MethodDelete,
		"/{uuid}",
		common.CreateHandler(
			h.DeleteGroup,
			h.DeleteGroupRequest,
			h.DeleteGroupResponse,
		),
	)

	router.Method(
		http.MethodDelete,
		"/{uuid}/members/{badgenumber}",
		common.CreateHandler(
			h.RemoveGroupMember,
			h.RemoveGroupMemberRequest,
			h.RemoveGroupMemberResponse,
		),
	)
}
