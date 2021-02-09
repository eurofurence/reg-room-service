package config

import (
	"github.com/eurofurence/reg-room-service/docs"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadConfiguration_noFilename(t *testing.T) {
	docs.Description("empty configuration filename is an error")
	err := LoadConfiguration("")
	require.NotNil(t, err)
	require.Equal(t, "no configuration filename provided", err.Error())
}
