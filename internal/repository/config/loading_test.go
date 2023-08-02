package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/eurofurence/reg-room-service/docs"
)

func TestLoadConfiguration_noFilename(t *testing.T) {
	docs.Description("empty configuration filename is an error")
	err := LoadConfiguration("")
	require.NotNil(t, err)
	require.Equal(t, "no configuration filename provided", err.Error())
}
