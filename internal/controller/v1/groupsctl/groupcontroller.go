package groupsctl

import (
	"github.com/eurofurence/reg-room-service/internal/application/web"
	"net/http"

	groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"

	"github.com/go-chi/chi/v5"
)

// Controller implements methods which satisfy the endpoint format
// in the `common` package.
type Controller struct {
	svc groupservice.Service
}

// InitRoutes creates the Controller instance and sets up all routes on it.
func InitRoutes(router chi.Router, svc groupservice.Service) {
	h := &Controller{
		svc: svc,
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
		web.CreateHandler(
			h.ListGroups,
			h.ListGroupsRequest,
			h.ListGroupsResponse,
		),
	)

	router.Method(
		http.MethodGet,
		"/my",
		web.CreateHandler(
			h.FindMyGroup,
			h.FindMyGroupRequest,
			h.FindMyGroupResponse,
		),
	)

	router.Method(
		http.MethodGet,
		"/{uuid}",
		web.CreateHandler(
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
		web.CreateHandler(
			h.CreateGroup,
			h.CreateGroupRequest,
			h.CreateGroupResponse,
		),
	)

	router.Method(
		http.MethodPost,
		"/{uuid}/members/{badgenumber}",
		web.CreateHandler(
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
		web.CreateHandler(
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
		web.CreateHandler(
			h.DeleteGroup,
			h.DeleteGroupRequest,
			h.DeleteGroupResponse,
		),
	)

	router.Method(
		http.MethodDelete,
		"/{uuid}/members/{badgenumber}",
		web.CreateHandler(
			h.RemoveGroupMember,
			h.RemoveGroupMemberRequest,
			h.RemoveGroupMemberResponse,
		),
	)
}
