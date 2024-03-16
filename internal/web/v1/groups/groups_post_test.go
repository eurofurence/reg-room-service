package groups

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/require"
// )

// // ======================= CreateGroup ==========================

// func TestCreateGroupRequest(t *testing.T) {
// 	testUUID := uuid.NewString()
// 	tests := []struct {
// 		name      string
// 		inputJSON []byte
// 		expected  *CreateGroupRequest
// 	}{
// 		{
// 			name: "Should successfully create request with provided parameters",
// 			inputJSON: marshalGroup(t, modelsv1.Group{
// 				ID:    testUUID,
// 				Name:  "test",
// 				Flags: []string{"test", "test2"},
// 				Owner: 10,
// 				Members: []modelsv1.Member{
// 					{
// 						ID:       10,
// 						Nickname: "Jumpy",
// 					},
// 					{
// 						ID:       11,
// 						Nickname: "Tabalon",
// 					},
// 				},
// 			}),
// 			expected: &CreateGroupRequest{
// 				Group: modelsv1.Group{
// 					ID:    testUUID,
// 					Name:  "test",
// 					Flags: []string{"test", "test2"},
// 					Owner: 10,
// 					Members: []modelsv1.Member{
// 						{
// 							ID:       10,
// 							Nickname: "Jumpy",
// 						},
// 						{
// 							ID:       11,
// 							Nickname: "Tabalon",
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name:      "Should fail to create request with invalid json",
// 			inputJSON: []byte(`{"id": "123`),
// 			expected:  nil,
// 		},
// 		{
// 			name:      "Should fail to create request with unknown field",
// 			inputJSON: []byte(`{"unknown": "field"}`),
// 			expected:  nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r, err := http.NewRequest(http.MethodPost, "http://test.test", bytes.NewBuffer(tt.inputJSON))
// 			require.NoError(t, err)
// 			req, err := (&Handler{}).CreateGroupRequest(r, httptest.NewRecorder())
// 			if tt.expected == nil {
// 				require.Error(t, err)
// 				require.Nil(t, req)
// 			} else {
// 				require.NoError(t, err)
// 				require.Equal(t, tt.expected, req)
// 			}
// 		})
// 	}
// }

// func marshalGroup(t *testing.T, g modelsv1.Group) []byte {
// 	t.Helper()

// 	b, err := json.Marshal(&g)
// 	require.NoError(t, err)

// 	return b
// }

// // ==============================================================
