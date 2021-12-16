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

func tstValidatePublicBookingCode(t *testing.T, value string, errMessage string) {
	errs := validationErrors{}
	config := goLiveConfig{Public: publicGoLiveConfig{BookingCode: value}}
	validateGoLiveConfiguration(errs, config)
	require.Equal(t, []string{errMessage}, errs["go_live.public.booking_code"])
}

func TestValidatePublicGoLiveConfiguration_codeEmpty(t *testing.T) {
	docs.Description("validation should catch an empty booking code")
	tstValidatePublicBookingCode(t, "", "value '' cannot be empty")
}

func tstValidatePublicStartIsoDatetime(t *testing.T, value string, errMessage string) {
	errs := validationErrors{}
	config := goLiveConfig{Public: publicGoLiveConfig{StartIsoDatetime: value}}
	validateGoLiveConfiguration(errs, config)
	require.Equal(t, []string{errMessage}, errs["go_live.public.start_iso_datetime"])
}

func TestValidatePublicGoLiveConfiguration_startTimeEmpty(t *testing.T) {
	docs.Description("validation should catch an empty go live time")
	tstValidatePublicStartIsoDatetime(t, "", "value '' cannot be empty")
}

func TestValidatePublicGoLiveConfiguration_startTimeInvalid(t *testing.T) {
	docs.Description("validation should catch an empty go live time")
	tstValidatePublicStartIsoDatetime(t, "2019-02-29T25:14:31-12:00", "value '2019-02-29T25:14:31-12:00' is not a valid go live time, format is 2006-01-02T15:04:05-07:00")
}

func tstValidateTokenPublicKey(t *testing.T, value string, errMessage string) {
	errs := validationErrors{}
	config := securityConfig{TokenPublicKeyPEM: value}
	validateSecurityConfiguration(errs, config)
	require.Equal(t, []string{errMessage}, errs["security.token_public_key_PEM"])
}

func TestValidateSecurityConfiguration_publicKeyEmpty(t *testing.T) {
	docs.Description("validation should catch an empty public key")
	tstValidateTokenPublicKey(t, "", "value '' cannot be empty")
}

func TestValidateSecurityConfiguration_publicKeyInvalid(t *testing.T) {
	docs.Description("validation should catch an invalid PEM public key")
	tstValidateTokenPublicKey(t, "MIIBIjANBgkqhkiG9", "value '(omitted)' is not a valid RSA256 PEM key")
}

func TestValidateSecurityConfiguration_publicKeyWrongKeySize(t *testing.T) {
	rsa128key := `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDfMWoMHiq6t5gVgcptidocUzc6
bED0dtVGFrYP/xD+Ew/Ecv37f1TXed2h6BFTf5luTB0DDWY7eolmhPsP1VFL8aSA
3uoH9IN6pJtEB/KZCSxxGjgTzGm0wgD/hDTtGZk+yricipoKMZW4TbS7kSfVj6JL
rMFjtxXoOTJyE+6t/QIDAQAB
-----END PUBLIC KEY-----`
	docs.Description("validation should catch a PEM public key of the wrong size")
	tstValidateTokenPublicKey(t, rsa128key, "value '(omitted)' has wrong key size, must be 256 (2048bit) or 512 (4096bit)")
}
