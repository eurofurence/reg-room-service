package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"

	"github.com/eurofurence/reg-room-service/docs"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/authservice"
)

// --- test setup ---

func TestMain(m *testing.M) {
	tstSetup()
	code := m.Run()
	os.Exit(code)
}

var authServiceMock authservice.Mock

var securityConfig256 = config.SecurityConfig{
	Fixed: config.FixedTokenConfig{
		API: valid_api_token,
	},
	Oidc: config.OpenIDConnectConfig{
		IDTokenCookieName:     "JWT",
		AccessTokenCookieName: "AUTH",
		TokenPublicKeysPEM:    []string{pemPublicKeyRS256},
		AdminGroup:            "admin",
	},
}

func tstSetup() {
	aulogging.SetupNoLoggerForTesting()
	authServiceMock = authservice.CreateMock()
	// call to make sure to fill the parsed pems for the tests
	CheckRequestAuthorization(&securityConfig256)
}

// --- tokens and security config ---

const valid_JWT_id_is_not_staff_sub101 = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdF9oYXNoIjoidDdWYkV5NVQ3STdYSlh3VHZ4S3hLdyIsImF1ZCI6WyIxNGQ5ZjM3YS0xZWVjLTQ3YzktYTk0OS01ZjFlYmRmOWM4ZTUiXSwiYXV0aF90aW1lIjoxNTE2MjM5MDIyLCJlbWFpbCI6ImpzcXVpcnJlbF9naXRodWJfOWE2ZEBwYWNrZXRsb3NzLmRlIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV4cCI6MjA3NTEyMDgxNiwiZ3JvdXBzIjpbInNvbWVncm91cCJdLCJpYXQiOjE1MTYyMzkwMjIsImlzcyI6Imh0dHA6Ly9pZGVudGl0eS5sb2NhbGhvc3QvIiwianRpIjoiNDA2YmUzZTQtZjRlOS00N2I3LWFjNWYtMDZiOTI3NDMyODQ4IiwibmFtZSI6IkpvaG4gRG9lIiwibm9uY2UiOiIzMGM4M2MxM2M5MTc5ODA0YWEwZjliMzkzNDI1OWQ3NSIsInJhdCI6MTY3NTExNzE3Nywic2lkIjoiZDdiOGZlN2EtMDc5YS00NTk2LThlNTMtYTYwZjg2YTA4YWM2Iiwic3ViIjoiMTAxIn0.ntHz3G7LLtHC3pJ1PoWJoG3mnzg96IIcP3LAV4V1CcKYMFoKVQfh7MiOdRXpiB-_j4QFE7O-za3mynwFqRbF3_Tw_Sp7Zsgk9OUPo2Mk3VBSl9yPIU4pmc8v7nrmaAVOQLyjglVG7NLRWLpx0oIG8SSN0d75PBI5iLyQ0H7Zu0npEu6xekHeAYAg9DHQxqZInzom72aLmHdtG7tOqOgN0XphiK7zmIqm5aCg7R9_J9s0UU0g16_Phxm3DaynufGCjEPE2YrSL7hY9UVT2nfrHO7MvVOEKMG3RaKUDjzqOkLawz9TcUJlUTBc1J-91zYbdXLHYT_2b4EW_qa1C-P3Ow`

const valid_JWT_id_is_staff_sub202 = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdF9oYXNoIjoidDdWYkV5NVQ3STdYSlh3VHZ4S3hLdyIsImF1ZCI6WyIxNGQ5ZjM3YS0xZWVjLTQ3YzktYTk0OS01ZjFlYmRmOWM4ZTUiXSwiYXV0aF90aW1lIjoxNTE2MjM5MDIyLCJlbWFpbCI6ImpzcXVpcnJlbF9naXRodWJfOWE2ZEBwYWNrZXRsb3NzLmRlIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV4cCI6MjA3NTEyMDgxNiwiZ3JvdXBzIjpbInN0YWZmIl0sImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LmxvY2FsaG9zdC8iLCJqdGkiOiI0MDZiZTNlNC1mNGU5LTQ3YjctYWM1Zi0wNmI5Mjc0MzI4NDgiLCJuYW1lIjoiSm9obiBTdGFmZiIsIm5vbmNlIjoiMzBjODNjMTNjOTE3OTgwNGFhMGY5YjM5MzQyNTlkNzUiLCJyYXQiOjE2NzUxMTcxNzcsInNpZCI6ImQ3YjhmZTdhLTA3OWEtNDU5Ni04ZTUzLWE2MGY4NmEwOGFjNiIsInN1YiI6IjIwMiJ9.pM-jMGdjwNvHQMov8JQpRa79CBjHAUHpwElYRvUz_DvhkqcG4SrntVruAlJRS8D9CccflKeTjSEfOiS2l52p0qQ7ZeNPSRQ9nsr_EHDpB7UqcUszqVaBWtIhwkiwca_sxe-8L9A9hPSe_kH9dhDHVbhUsj9vp0HBIV89mtH3i3D6s3quRYtCe9puepkmyf5JC-2TSGoSiyURoFdqXSNRPEuv1FhlyVICqylfkINceCe8dt7lJCrhOc8R-11vA-SRsrBhdxBvYjT29hhQO3eHgJenPufjJPj6kYSWvN91U3KcsffoBmu-C1A7zBLu-zmWBHh4RkYWqbZpNr59TIpRSw`

const valid_JWT_id_is_admin_sub1234567890 = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdF9oYXNoIjoidDdWYkV5NVQ3STdYSlh3VHZ4S3hLdyIsImF1ZCI6WyIxNGQ5ZjM3YS0xZWVjLTQ3YzktYTk0OS01ZjFlYmRmOWM4ZTUiXSwiYXV0aF90aW1lIjoxNTE2MjM5MDIyLCJlbWFpbCI6ImpzcXVpcnJlbF9naXRodWJfOWE2ZEBwYWNrZXRsb3NzLmRlIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV4cCI6MjA3NTEyMDgxNiwiZ3JvdXBzIjpbInN0YWZmIiwiYWRtaW4iXSwiaWF0IjoxNTE2MjM5MDIyLCJpc3MiOiJodHRwOi8vaWRlbnRpdHkubG9jYWxob3N0LyIsImp0aSI6IjQwNmJlM2U0LWY0ZTktNDdiNy1hYzVmLTA2YjkyNzQzMjg0OCIsIm5hbWUiOiJKb2huIEFkbWluIiwibm9uY2UiOiIzMGM4M2MxM2M5MTc5ODA0YWEwZjliMzkzNDI1OWQ3NSIsInJhdCI6MTY3NTExNzE3Nywic2lkIjoiZDdiOGZlN2EtMDc5YS00NTk2LThlNTMtYTYwZjg2YTA4YWM2Iiwic3ViIjoiMTIzNDU2Nzg5MCJ9.DRKPy0Rq-r5-Va6W5coT91JpDV2RkhYjniqIJmmPzOq3LphzRrlDKioDns4ilMxMEpfhFcmv87yOdPsPijUhEqy1a93BeJYMyU7DMdQBtD8R9oYU_-FmqS5dM9ZrBCZZUxTBeNBl2JGI-H1_IIqUH65PodoijO4N5ayw43q5xT1KP7PJKZ9YiMSsa4tUOp0R_Ay51DTIuti21TqqbSCC66sGH_1e1eeuhwBoU7Iws4oeepTRZ_XOdpn_YzTViPs7Byua-zohYgQYthDoCvLVfJOr77BV2vTUrQZfRca7prizXbVuQyxQJEpOBIuke29Gye6Qfbwpb4rMaza3fTLhZg`

const invalid_JWT_id_is_admin_wrong_signature = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdF9oYXNoIjoidDdWYkV5NVQ3STdYSlh3VHZ4S3hLdyIsImF1ZCI6WyIxNGQ5ZjM3YS0xZWVjLTQ3YzktYTk0OS01ZjFlYmRmOWM4ZTUiXSwiYXV0aF90aW1lIjoxNTE2MjM5MDIyLCJlbWFpbCI6ImpzcXVpcnJlbF9naXRodWJfOWE2ZEBwYWNrZXRsb3NzLmRlIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV4cCI6MjA3NTEyMDgxNiwiZ3JvdXBzIjpbInN0YWZmIiwiYWRtaW4iXSwiaWF0IjoxNTE2MjM5MDIyLCJpc3MiOiJodHRwOi8vaWRlbnRpdHkubG9jYWxob3N0LyIsImp0aSI6IjQwNmJlM2U0LWY0ZTktNDdiNy1hYzVmLTA2YjkyNzQzMjg0OCIsIm5hbWUiOiJKb2huIEFkbWluIiwibm9uY2UiOiIzMGM4M2MxM2M5MTc5ODA0YWEwZjliMzkzNDI1OWQ3NSIsInJhdCI6MTY3NTExNzE3Nywic2lkIjoiZDdiOGZlN2EtMDc5YS00NTk2LThlNTMtYTYwZjg2YTA4YWM2Iiwic3ViIjoiMTIzNDU2Nzg5MCJ9.DRKPy0Rq-r5-Va6W5coT91JpDV2RkhYjniqIJmmPzOq3LphzRrlDKioDns4ilMxMEpfhFcmv87yOdPsPijUhEqy1a93BeJYMyU7DMdQBtD8R9oYU_-FmqS5dM9ZrBCZZUxTBeNBl2JGI-H1_IIqUH65PodoijO4N5ayw43q5xT1KP7PJKZ9YiMSsa4tUOp0R_Ay51DTIuti21TqqbSCC66sGH_1e1eeuhwBoU7Iws4oeepTRZ_XOdpn_YzTViPs7Byua-zohYgQYthDoCvLVfJOr77BV2vTUrQZfRca7prizXbVuQyxQJEpOBIuke29Gye6Qfbwpb4r3fTLhZg`

const expired_valid_JWT_id_is_admin = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdF9oYXNoIjoidDdWYkV5NVQ3STdYSlh3VHZ4S3hLdyIsImF1ZCI6WyIxNGQ5ZjM3YS0xZWVjLTQ3YzktYTk0OS01ZjFlYmRmOWM4ZTUiXSwiYXV0aF90aW1lIjoxNTE2MjM5MDIyLCJlbWFpbCI6ImpzcXVpcnJlbF9naXRodWJfOWE2ZEBwYWNrZXRsb3NzLmRlIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV4cCI6MTUxNjIzOTAyMywiZ3JvdXBzIjpbInN0YWZmIiwiYWRtaW4iXSwiaWF0IjoxNTE2MjM5MDIyLCJpc3MiOiJodHRwOi8vaWRlbnRpdHkubG9jYWxob3N0LyIsImp0aSI6IjQwNmJlM2U0LWY0ZTktNDdiNy1hYzVmLTA2YjkyNzQzMjg0OCIsIm5hbWUiOiJKb2huIEFkbWluIiwibm9uY2UiOiIzMGM4M2MxM2M5MTc5ODA0YWEwZjliMzkzNDI1OWQ3NSIsInJhdCI6MTY3NTExNzE3Nywic2lkIjoiZDdiOGZlN2EtMDc5YS00NTk2LThlNTMtYTYwZjg2YTA4YWM2Iiwic3ViIjoiMTIzNDU2Nzg5MCJ9.BM3c4PccnY7AuazPgk2eBq8_vBO5iAEqffjc8NosxJTeVWaZFRL8Zz7WWrZ4EVPAWTwy5AlvR3Vva6n82VqDQMA0nCQTrTqv73aFvn6b2A81cpyxVGUZdIOAPR7mS0dozzMaR2H9rZ3t946ppqTYHG2GiovJGuABS8AtabG0_dlKCFT4-910ndYIwtH8V71FtplvHLZzwg1X7d5EP8Mhwp_iYwMbu3BubRwHdzKByH6NXEkEbRPTt9RIBf4xZWcXXj4oKiW2U3h5p0wUwfwWbrYHU8STHmpcOJJyl9oDrgWtrZFyGqNKw3elrIoWKRFMRhVbE2WVQU_yMOnNV6QQYw`

const invalid_api_token = "invalid"

const valid_api_token = "api-token-for-testing-must-be-pretty-long"

const invalid_access_token = "invalid-access"

const valid_access_token = "valid-access"

// --- test case helpers ---

func tstRequire(t *testing.T, actualMsg string, actualErr error, expectedMsg string, expectedErr string) {
	if expectedErr != "" {
		require.NotNil(t, actualErr)
		require.Contains(t, actualErr.Error(), expectedErr)
	} else {
		require.Nil(t, actualErr)
	}
	require.Equal(t, expectedMsg, actualMsg)
}

func tstRequireNoAuthServiceCall(t *testing.T) {
	require.Equal(t, len(authServiceMock.Recording()), 0)
}

func tstRequireAuthServiceCall(t *testing.T, idToken string, accToken string) {
	recording := authServiceMock.Recording()
	require.Equal(t, len(recording), 1)
	require.Equal(t, recording[0], fmt.Sprintf("userinfo %s %s", idToken, accToken))
}

func tstNothingTestCase(t *testing.T, expectedMsg string, expectedErr string) context.Context {
	ctx, actualMsg, actualErr := checkAllAuthentication(context.Background(), http.MethodGet, "/not/health", &securityConfig256, "", "", "", "")
	tstRequire(t, actualMsg, actualErr, expectedMsg, expectedErr)
	return ctx
}

func tstApiTokenTestCase(t *testing.T, apiTokenHeaderValue string, expectedUserMsg string, expectedLoggedErr string) context.Context {
	ctx, actualMsg, actualErr := checkAllAuthentication(context.Background(), http.MethodGet, "/not/health", &securityConfig256, apiTokenHeaderValue, "", "", "")
	tstRequire(t, actualMsg, actualErr, expectedUserMsg, expectedLoggedErr)
	return ctx
}

func tstAuthHeaderTestCase(t *testing.T, authHeaderValue string, expectedUserMsg string, expectedLoggedErr string) context.Context {
	ctx, actualMsg, actualErr := checkAllAuthentication(context.Background(), http.MethodGet, "/not/health", &securityConfig256, "", authHeaderValue, "", "")
	tstRequire(t, actualMsg, actualErr, expectedUserMsg, expectedLoggedErr)
	return ctx
}

func tstCookiesTestCase(t *testing.T, idTokenCookieValue string, accessTokenCookieValue string, expectedUserMsg string, expectedLoggedErr string) context.Context {
	ctx, actualMsg, actualErr := checkAllAuthentication(context.Background(), http.MethodGet, "/not/health", &securityConfig256, "", "", idTokenCookieValue, accessTokenCookieValue)
	tstRequire(t, actualMsg, actualErr, expectedUserMsg, expectedLoggedErr)
	return ctx
}

// --- test cases ---

func TestHealthCheck(t *testing.T) {
	docs.Description("health check is allowed through")
	ctx, actualMsg, actualErr := checkAllAuthentication(context.Background(), http.MethodGet, "/", &securityConfig256, "", "", "", "")
	tstRequire(t, actualMsg, actualErr, "", "")
	require.Nil(t, ctx.Value(common.CtxKeyAPIKey{}))
	require.Nil(t, ctx.Value(common.CtxKeyIDToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyAccessToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyClaims{}))
}

func TestNothingProvided(t *testing.T) {
	docs.Description("not providing any authorization fails for this service")
	ctx := tstNothingTestCase(t, "you must be logged in for this operation", "no authorization presented")
	require.Nil(t, ctx.Value(common.CtxKeyAPIKey{}))
	require.Nil(t, ctx.Value(common.CtxKeyIDToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyAccessToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyClaims{}))
}

func TestApiTokenInvalid(t *testing.T) {
	docs.Description("invalid Api Token values are rejected")
	tstApiTokenTestCase(t, invalid_api_token, "invalid api token", "token doesn't match the configured value")
}

func TestApiTokenValid(t *testing.T) {
	docs.Description("valid Api Token values authorize as api user")
	ctx := tstApiTokenTestCase(t, valid_api_token, "", "")
	require.Equal(t, valid_api_token, ctx.Value(common.CtxKeyAPIKey{}))
	require.Nil(t, ctx.Value(common.CtxKeyIDToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyAccessToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyClaims{}))
}

func TestAccessTokenAuthDisabled(t *testing.T) {
	docs.Description("any access token is rejected if no userinfo endpoint is available")
	authServiceMock.Reset()
	ctx := tstAuthHeaderTestCase(t, valid_access_token, "invalid bearer token", "request failed access token check, denying: no userinfo endpoint configured")
	require.Nil(t, ctx.Value(common.CtxKeyAPIKey{}))
	require.Nil(t, ctx.Value(common.CtxKeyIDToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyAccessToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyClaims{}))
	tstRequireNoAuthServiceCall(t)
}

func TestAccessTokenInvalid(t *testing.T) {
	docs.Description("invalid access token is rejected after rejected call to userinfo")
	authServiceMock.Reset()
	authServiceMock.Enable()
	authServiceMock.SimulateGetError(authservice.UnauthorizedError)
	ctx := tstAuthHeaderTestCase(t, invalid_access_token, "invalid bearer token", "request failed access token check, denying: got unauthorized from userinfo endpoint")
	require.Nil(t, ctx.Value(common.CtxKeyAPIKey{}))
	require.Nil(t, ctx.Value(common.CtxKeyIDToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyAccessToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyClaims{}))
	tstRequireAuthServiceCall(t, "", invalid_access_token)
}

func TestAccessTokenUserinfoUnavailable(t *testing.T) {
	docs.Description("valid access token is rejected when auth service unresponsive")
	authServiceMock.Reset()
	authServiceMock.Enable()
	authServiceMock.SimulateGetError(authservice.DownstreamError)
	ctx := tstAuthHeaderTestCase(t, valid_access_token, "invalid bearer token", "request failed access token check, denying: downstream unavailable - see log for details")
	require.Nil(t, ctx.Value(common.CtxKeyAPIKey{}))
	require.Nil(t, ctx.Value(common.CtxKeyIDToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyAccessToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyClaims{}))
	tstRequireAuthServiceCall(t, "", valid_access_token)
}

func TestAccessTokenValid(t *testing.T) {
	docs.Description("valid access token is accepted after call to userinfo")
	authServiceMock.Reset()
	authServiceMock.Enable()
	authServiceMock.SetupResponse("", valid_access_token, authservice.UserInfoResponse{
		Email:         "jsquirrel_github_9a6d@packetloss.de",
		EmailVerified: true,
		Groups:        []string{"admin"},
		Subject:       "subject",
		Name:          "name",
	})
	ctx := tstAuthHeaderTestCase(t, valid_access_token, "", "")
	require.Nil(t, ctx.Value(common.CtxKeyAPIKey{}))
	require.Nil(t, ctx.Value(common.CtxKeyIDToken{}))
	require.Equal(t, valid_access_token, ctx.Value(common.CtxKeyAccessToken{}))
	actualClaims, _ := ctx.Value(common.CtxKeyClaims{}).(*common.AllClaims)
	require.NotNil(t, actualClaims)
	expectedClaims := common.AllClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "",
			Subject:  "subject",
			Audience: jwt.ClaimStrings{""},
		},
		CustomClaims: common.CustomClaims{
			EMail:         "jsquirrel_github_9a6d@packetloss.de",
			EMailVerified: true,
			Groups:        []string{"admin"},
			Name:          "name",
		},
	}
	require.EqualValues(t, expectedClaims, *actualClaims)
	tstRequireAuthServiceCall(t, "", valid_access_token)
}

func TestCookiesValidWithUserinfo(t *testing.T) {
	docs.Description("valid cookies are accepted if matching results of userinfo lookup")
	authServiceMock.Reset()
	authServiceMock.Enable()
	authServiceMock.SetupResponse(valid_JWT_id_is_admin_sub1234567890, valid_access_token, authservice.UserInfoResponse{
		Subject:       "1234567890",
		Name:          "John Admin",
		Email:         "jsquirrel_github_9a6d@packetloss.de",
		EmailVerified: true,
		Groups:        []string{"admin"},
	})
	ctx := tstCookiesTestCase(t, valid_JWT_id_is_admin_sub1234567890, valid_access_token, "", "")
	require.Nil(t, ctx.Value(common.CtxKeyAPIKey{}))
	require.Equal(t, valid_JWT_id_is_admin_sub1234567890, ctx.Value(common.CtxKeyIDToken{}))
	require.Equal(t, valid_access_token, ctx.Value(common.CtxKeyAccessToken{}))
	actualClaims, _ := ctx.Value(common.CtxKeyClaims{}).(*common.AllClaims)
	require.NotNil(t, actualClaims)
	expectedClaims := common.AllClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "",
			Subject:  "1234567890",
			Audience: jwt.ClaimStrings{""},
		},
		CustomClaims: common.CustomClaims{
			EMail:         "jsquirrel_github_9a6d@packetloss.de",
			EMailVerified: true,
			Groups:        []string{"admin"},
			Name:          "John Admin",
		},
	}
	require.EqualValues(t, expectedClaims, *actualClaims)
	tstRequireAuthServiceCall(t, valid_JWT_id_is_admin_sub1234567890, valid_access_token)
}

func TestCookiesInvalidJwt(t *testing.T) {
	docs.Description("invalid cookies (jwt signature broken) are rejected even before userinfo call")
	authServiceMock.Reset()
	authServiceMock.Enable()
	ctx := tstCookiesTestCase(t, invalid_JWT_id_is_admin_wrong_signature, valid_access_token, "invalid id token in cookie", "crypto/rsa: verification error")
	require.Nil(t, ctx.Value(common.CtxKeyAPIKey{}))
	require.Nil(t, ctx.Value(common.CtxKeyIDToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyAccessToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyClaims{}))
	tstRequireNoAuthServiceCall(t)
}

func TestCookiesExpiredJwt(t *testing.T) {
	docs.Description("invalid cookies (jwt expired) are rejected even before userinfo call")
	authServiceMock.Reset()
	authServiceMock.Enable()
	ctx := tstCookiesTestCase(t, expired_valid_JWT_id_is_admin, valid_access_token, "invalid id token in cookie", "token is expired by ")
	require.Nil(t, ctx.Value(common.CtxKeyAPIKey{}))
	require.Nil(t, ctx.Value(common.CtxKeyIDToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyAccessToken{}))
	require.Nil(t, ctx.Value(common.CtxKeyClaims{}))
	tstRequireNoAuthServiceCall(t)
}
