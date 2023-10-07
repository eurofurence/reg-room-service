package logging

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRequestID(t *testing.T) {
	tests := []struct {
		name         string
		inputContext context.Context
		expectedID   string
	}{
		{
			name:         "Should return valid requestID",
			inputContext: context.WithValue(context.Background(), RequestIDKey, "valid"),
			expectedID:   "valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedID, GetRequestID(tt.inputContext))
		})
	}
}
