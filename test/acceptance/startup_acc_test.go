package acceptance

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/eurofurence/reg-room-service/docs"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
)

// ----------------------------------------------------------
// acceptance tests for conditions that abort service startup
// ----------------------------------------------------------

func TestMissingConfiguration(t *testing.T) {
	docs.Given("given a missing configuration file")
	configFile := "../resources/i-am-missing.yaml"

	docs.When("when the service is started")
	_, err := config.UnmarshalFromYamlConfiguration(configFile)

	docs.Then("then it aborts with a useful error message")
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "i-am-missing.yaml: ")
}

func TestConfigurationSyntaxInvalid(t *testing.T) {
	docs.Given("given a syntactically invalid configuration file")
	configFile := "../resources/config-syntaxerror.yaml"

	docs.When("when the service is started")
	_, err := config.UnmarshalFromYamlConfiguration(configFile)

	docs.Then("then it aborts with a useful error message")
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "yaml: unmarshal errors:\n  line 7: ")
}

func TestConfigurationValidationFailure(t *testing.T) {
	docs.Given("given an invalid configuration file with valid syntax")
	configFile := "../resources/config-dataerror.yaml"

	docs.When("when the service is started")
	conf, err := config.UnmarshalFromYamlConfiguration(configFile)
	require.Nil(t, err)
	conf.AddDefaults()

	docs.Then("then it aborts with a useful error message")
	err = conf.Validate()
	require.NotNil(t, err)
	require.Equal(t, "configuration validation error, see log output for details", err.Error())
}

func TestConfigurationDefaults(t *testing.T) {
	docs.Given("given a minimal configuration file")
	configFile := "../resources/config-minimal.yaml"

	docs.When("when the service is started")
	conf, err := config.UnmarshalFromYamlConfiguration(configFile)
	require.Nil(t, err)
	conf.AddDefaults()

	docs.Then("then loading the configuration is successful and defaults are set")
	err = conf.Validate()
	require.Nil(t, err)
	require.Equal(t, ":8081", net.JoinHostPort(conf.Server.BaseAddress, fmt.Sprintf("%d", conf.Server.Port)))
}

func TestConfigurationFull(t *testing.T) {
	docs.Given("given a maximal configuration file")
	configFile := "../resources/config-maximal.yaml"

	docs.When("when the service is started")
	conf, err := config.UnmarshalFromYamlConfiguration(configFile)
	require.Nil(t, err)
	conf.AddDefaults()

	docs.Then("then loading the configuration is successful and all values are set")
	err = conf.Validate()
	require.Nil(t, err)
	require.Equal(t, ":12345", net.JoinHostPort(conf.Server.BaseAddress, fmt.Sprintf("%d", conf.Server.Port)))
	require.Equal(t, "Linköping", conf.GoLive.Public.BookingCode)
	// require.Equal(t, time.Date(2020, 11, 6, 21, 22, 23, 0, time.UTC).Unix(), config.PublicBookingStartTime().Unix())
}
