package groups

//
//import (
//	"context"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//	"github.com/stretchr/testify/suite"
//
//	"github.com/eurofurence/reg-room-service/internal/repository/database/inmemorydb"
//	groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"
//	"github.com/eurofurence/reg-room-service/internal/web/v1/util"
//)
//
//type groupsSuite struct {
//	suite.Suite
//	*require.Assertions
//
//	ctx  context.Context
//	ctrl *Controller
//}
//
//func TestGroupsSuite(t *testing.T) {
//	g := new(groupsSuite)
//	suite.Run(t, g)
//}
//
//func (g *groupsSuite) SetupTest() {
//	g.Assertions = require.New(g.T())
//	g.ctx = context.Background()
//	g.ctrl = &Controller{groupservice.NewService(inmemorydb.New())}
//}
//
//func TestParseGroupMemberIDs(t *testing.T) {
//	tests := []struct {
//		name        string
//		input       string
//		expected    []string
//		expectedErr bool
//	}{
//		{
//			name:        "Test with valid input",
//			input:       "1,2,3,4,5,6,7",
//			expected:    []string{"1", "2", "3", "4", "5", "6", "7"},
//			expectedErr: false,
//		},
//		{
//			name:        "Test with invalid input",
//			input:       "1,2,3,4,5,6,7,abc123",
//			expected:    nil,
//			expectedErr: true,
//		},
//		{
//			name:        "Test with brackets",
//			input:       "[{1,2,3,4,5,6,7]}]]",
//			expected:    nil,
//			expectedErr: true,
//		},
//	}
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			memberIDs, err := util.ParseMemberIDs(tc.input)
//			require.Equal(t, tc.expected, memberIDs)
//			if tc.expectedErr {
//				require.Error(t, err)
//			} else {
//				require.NoError(t, err)
//			}
//		})
//	}
//}
