package groups

// import (
// 	"context"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/stretchr/testify/require"
// )

// // ======================= ListGroups ==========================

// func TestListGroupsRequest(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		inputURL string
// 		expected *ListGroupsRequest
// 	}{

// 		{
// 			name:     "Should successfully create ListGroupRequest",
// 			inputURL: "http://test.groups.list/groups?member_ids=1,2,3,4,5&min_size=3&max_size=10",
// 			expected: &ListGroupsRequest{
// 				MemberIDs: []string{"1", "2", "3", "4", "5"},
// 				MinSize:   3,
// 				MaxSize:   10,
// 			},
// 		},
// 		{
// 			name:     "Should return error when member ID is non numeric",
// 			inputURL: "http://test.groups.list/groups?member_ids=1,2,E&min_size=3",
// 			expected: nil,
// 		},
// 		{
// 			name:     "Should return error on negative min size",
// 			inputURL: "http://test.groups.list/groups?member_ids=1,2,3&min_size=-10",
// 			expected: nil,
// 		},
// 		{
// 			name:     "Should return error on negative max size",
// 			inputURL: "http://test.groups.list/groups?member_ids=1,2,3&&max_size=-10",
// 			expected: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := &Controller{}
// 			r, err := http.NewRequest(http.MethodGet, tt.inputURL, nil)
// 			require.NoError(t, err)
// 			req, err := h.ListGroupsRequest(r, httptest.NewRecorder())

// 			if tt.expected != nil {
// 				require.NoError(t, err)
// 				require.Equal(t, tt.expected, req)
// 			} else {
// 				require.Nil(t, req)
// 				require.Error(t, err)
// 			}
// 		})
// 	}
// }

// // =============================================================

// // ======================= FindGroupsByID ======================

// func TestFindGroupsByIDRequest(t *testing.T) {

// 	tests := []struct {
// 		name     string
// 		inputURL string
// 		expected *FindGroupByIDRequest
// 	}{
// 		{
// 			name:     "Should successfully parse uuid and return request",
// 			inputURL: "http://test.groups.list/groups/2868913d-4bef-477f-bf00-d2ee246caa3b",
// 			expected: &FindGroupByIDRequest{
// 				GroupID: "2868913d-4bef-477f-bf00-d2ee246caa3b",
// 			},
// 		},
// 		{
// 			name:     "Should fail if uuid is not valid",
// 			inputURL: "http://test.groups.list/groups/invalid",
// 			expected: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			uuid := tt.inputURL[strings.LastIndex(tt.inputURL, "/")+1:]

// 			rctx := chi.NewRouteContext()
// 			rctx.URLParams.Add("uuid", uuid)

// 			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)
// 			r, err := http.NewRequestWithContext(ctx, http.MethodGet, tt.inputURL, nil)
// 			require.NoError(t, err)

// 			res, err := (&Controller{}).FindGroupByIDRequest(r, httptest.NewRecorder())
// 			if tt.expected == nil {
// 				require.Error(t, err)
// 				require.Nil(t, res)
// 			} else {
// 				require.NoError(t, err)
// 				require.Equal(t, tt.expected, res)
// 			}
// 		})
// 	}
// }

// // =============================================================
