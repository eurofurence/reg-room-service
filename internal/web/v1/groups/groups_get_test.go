package groups

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListGroupsRequest(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected *ListGroupsRequest
	}{

		{
			name:     "Should successfully create ListGroupRequest",
			inputURL: "http://test.groups.list?member_ids=1,2,3,4,5&min_size=3&max_size=10",
			expected: &ListGroupsRequest{
				MemberIDs: []string{"1", "2", "3", "4", "5"},
				MinSize:   3,
				MaxSize:   10,
			},
		},
		{
			name:     "Should return error when member ID is non numeric",
			inputURL: "http://test.groups.list?member_ids=1,2,E&min_size=3",
			expected: nil,
		},
		{
			name:     "Should return error on negative min size",
			inputURL: "http://test.groups.list?member_ids=1,2,3&min_size=-10",
			expected: nil,
		},
		{
			name:     "Should return error on negative max size",
			inputURL: "http://test.groups.list?member_ids=1,2,3&&max_size=-10",
			expected: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}
			r, err := http.NewRequest(http.MethodGet, tt.inputURL, nil)
			require.NoError(t, err)
			req, err := h.ListGroupsRequest(r, httptest.NewRecorder())

			if tt.expected != nil {
				require.NoError(t, err)
				require.Equal(t, tt.expected, req)
			} else {
				require.Nil(t, req)
				require.Error(t, err)
			}
		})
	}
}
