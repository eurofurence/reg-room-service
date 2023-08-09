package groups

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseGroupMemberIDs(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []string
		expectedErr bool
	}{
		{
			name:        "Test with valid input",
			input:       "1,2,3,4,5,6,7",
			expected:    []string{"1", "2", "3", "4", "5", "6", "7"},
			expectedErr: false,
		},
		{
			name:        "Test with invalid input",
			input:       "1,2,3,4,5,6,7,abc123",
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "Test with brackets",
			input:       "[{1,2,3,4,5,6,7]}]]",
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			memberIDs, err := parseGroupMemberIDs(tc.input)
			require.Equal(t, tc.expected, memberIDs)
			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
