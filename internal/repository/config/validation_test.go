package config

import (
	"github.com/eurofurence/reg-room-service/docs"
	"github.com/stretchr/testify/require"
	"testing"
)

func tstValidatePort(t *testing.T, value string, errMessage string) {
	errs := validationErrors{}
	config := serverConfig{Port: value}
	validateServerConfiguration(errs, config)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{errMessage}, errs["server.port"])
}

func TestValidateServerConfiguration_empty(t *testing.T) {
	docs.Description("validation should catch an empty port configuration")
	tstValidatePort(t, "", "value '' cannot be empty")
}

func TestValidateServerConfiguration_numeric(t *testing.T) {
	docs.Description("validation should catch a non-numeric port configuration")
	tstValidatePort(t, "katze", "value 'katze' is not a valid port number")
}

func TestValidateServerConfiguration_tooHigh(t *testing.T) {
	docs.Description("validation should catch a port configuration that is out of range")
	tstValidatePort(t, "65536", "value '65536' is not a valid port number")
}

func TestValidateServerConfiguration_privileged(t *testing.T) {
	docs.Description("validation should not allow privileged ports")
	tstValidatePort(t, "1023", "value '1023' must be a nonprivileged port")
}

func tstValidateBookingCode(t *testing.T, value string, errMessage string) {
	errs := validationErrors{}
	config := goLiveConfig{BookingCode: value}
	validateGoLiveConfiguration(errs, config)
	require.Equal(t, []string{errMessage}, errs["go_live.booking_code"])
}

func TestValidateGoLiveConfiguration_codeEmpty(t *testing.T) {
	docs.Description("validation should catch an empty booking code")
	tstValidateBookingCode(t, "", "value '' cannot be empty")
}

func tstValidateStartIsoDatetime(t *testing.T, value string, errMessage string) {
	errs := validationErrors{}
	config := goLiveConfig{StartIsoDatetime: value}
	validateGoLiveConfiguration(errs, config)
	require.Equal(t, []string{errMessage}, errs["go_live.start_iso_datetime"])
}

func TestValidateGoLiveConfiguration_startTimeEmpty(t *testing.T) {
	docs.Description("validation should catch an empty go live time")
	tstValidateStartIsoDatetime(t, "", "value '' cannot be empty")
}

func TestValidateGoLiveConfiguration_startTimeInvalid(t *testing.T) {
	docs.Description("validation should catch an empty go live time")
	tstValidateStartIsoDatetime(t, "2019-02-29T25:14:31-12:00", "value '2019-02-29T25:14:31-12:00' is not a valid go live time, format is 2006-01-02T15:04:05-07:00")
}
