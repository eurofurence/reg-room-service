package jwt

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnzyis1ZjfNB0bBgKFMSv
vkTtwlvBsaJq7S5wA+kzeVOVpVWwkWdVha4s38XM/pa/yr47av7+z3VTmvDRyAHc
aT92whREFpLv9cj5lTeJSibyr/Mrm/YtjCZVWgaOYIhwrXwKLqPr/11inWsAkfIy
tvHWTxZYEcXLgAXFuUuaS3uF9gEiNQwzGTU1v0FqkqTBr4B8nW3HCN47XUu0t8Y0
e+lf4s4OxQawWD79J9/5d3Ry0vbV3Am1FtGJiJvOwRsIfVChDpYStTcHTCMqtvWb
V6L11BWkpzGXSW4Hv43qa+GSYOD2QU68Mb59oSk2OB+BtOLpJofmbGEGgvmwyCI9
MwIDAQAB
-----END PUBLIC KEY-----`

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAnzyis1ZjfNB0bBgKFMSvvkTtwlvBsaJq7S5wA+kzeVOVpVWw
kWdVha4s38XM/pa/yr47av7+z3VTmvDRyAHcaT92whREFpLv9cj5lTeJSibyr/Mr
m/YtjCZVWgaOYIhwrXwKLqPr/11inWsAkfIytvHWTxZYEcXLgAXFuUuaS3uF9gEi
NQwzGTU1v0FqkqTBr4B8nW3HCN47XUu0t8Y0e+lf4s4OxQawWD79J9/5d3Ry0vbV
3Am1FtGJiJvOwRsIfVChDpYStTcHTCMqtvWbV6L11BWkpzGXSW4Hv43qa+GSYOD2
QU68Mb59oSk2OB+BtOLpJofmbGEGgvmwyCI9MwIDAQABAoIBACiARq2wkltjtcjs
kFvZ7w1JAORHbEufEO1Eu27zOIlqbgyAcAl7q+/1bip4Z/x1IVES84/yTaM8p0go
amMhvgry/mS8vNi1BN2SAZEnb/7xSxbflb70bX9RHLJqKnp5GZe2jexw+wyXlwaM
+bclUCrh9e1ltH7IvUrRrQnFJfh+is1fRon9Co9Li0GwoN0x0byrrngU8Ak3Y6D9
D8GjQA4Elm94ST3izJv8iCOLSDBmzsPsXfcCUZfmTfZ5DbUDMbMxRnSo3nQeoKGC
0Lj9FkWcfmLcpGlSXTO+Ww1L7EGq+PT3NtRae1FZPwjddQ1/4V905kyQFLamAA5Y
lSpE2wkCgYEAy1OPLQcZt4NQnQzPz2SBJqQN2P5u3vXl+zNVKP8w4eBv0vWuJJF+
hkGNnSxXQrTkvDOIUddSKOzHHgSg4nY6K02ecyT0PPm/UZvtRpWrnBjcEVtHEJNp
bU9pLD5iZ0J9sbzPU/LxPmuAP2Bs8JmTn6aFRspFrP7W0s1Nmk2jsm0CgYEAyH0X
+jpoqxj4efZfkUrg5GbSEhf+dZglf0tTOA5bVg8IYwtmNk/pniLG/zI7c+GlTc9B
BwfMr59EzBq/eFMI7+LgXaVUsM/sS4Ry+yeK6SJx/otIMWtDfqxsLD8CPMCRvecC
2Pip4uSgrl0MOebl9XKp57GoaUWRWRHqwV4Y6h8CgYAZhI4mh4qZtnhKjY4TKDjx
QYufXSdLAi9v3FxmvchDwOgn4L+PRVdMwDNms2bsL0m5uPn104EzM6w1vzz1zwKz
5pTpPI0OjgWN13Tq8+PKvm/4Ga2MjgOgPWQkslulO/oMcXbPwWC3hcRdr9tcQtn9
Imf9n2spL/6EDFId+Hp/7QKBgAqlWdiXsWckdE1Fn91/NGHsc8syKvjjk1onDcw0
NvVi5vcba9oGdElJX3e9mxqUKMrw7msJJv1MX8LWyMQC5L6YNYHDfbPF1q5L4i8j
8mRex97UVokJQRRA452V2vCO6S5ETgpnad36de3MUxHgCOX3qL382Qx9/THVmbma
3YfRAoGAUxL/Eu5yvMK8SAt/dJK6FedngcM3JEFNplmtLYVLWhkIlNRGDwkg3I5K
y18Ae9n7dHVueyslrb6weq7dTkYDi3iOYRW8HRkIQh06wEdbxt0shTzAJvvCQfrB
jg/3747WSsf/zBTcHihTRBdAv6OmdhV4/dD5YBfLAkLrd+mX7iE=
-----END RSA PRIVATE KEY-----`

const valid_JWT_is_staff = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZ2xvYmFsIjp7Im5hbWUiOiJKb2huIERvZSIsInJvbGVzIjpbInN0YWZmIl19LCJpYXQiOjE1MTYyMzkwMjJ9.XKf8Eqrs7JmGSNhVUS-5RLIMnSuQHh65VMiUJPHaE5AEFZ55EY7MxsD2Sqdc6QV9cX0zA5weGXX2cGOAR0CNcjOGGsQSVogAcoEuwjve4WXLVvHPb41p95Jkbe9Md5bSPrk9oJwopJCVDI5DU1rLg0FIbt2yWORinZQiGvxZlPSZyNQuFoAXJXQPv4TNfTaBZcKhzUeO0u6_AQzIKGrF8VmbE4cMHq0fEAflnzroDmo9oJ-8dKJc2BNEyFQYHi9Jp3h3C85BvxEsdRzL3e9Qjw2SpFS0A8pPr4HEQikIn2nOEXav2RAcZMGN3YmdUeUBHwnfQ9ubY-0KilK9zNfGBw`

const valid_JWT_is_not_staff = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZ2xvYmFsIjp7Im5hbWUiOiJKb2huIERvZSIsInJvbGVzIjpbXX0sImlhdCI6MTUxNjIzOTAyMn0.IH3Q46k85RZsvgWD3wC9kNCtCRujTEOpzzCw6rqrKF4QoDcmn6Pd-Y2qQ8IZydrtzGrCu7yUiVziL634gxDlRvVliyHU6KkIMMsXDtnJWOGrKkpJgr_PZCA2LIlYD0GsXYzzQBuOg3eeXgidkGD7WVjHuKcuJe5By9nc6cTHlBHV-XeRIeCCy9jq10pbqyNv1kfjhdKuUQpFogV2JIKlTi3cR5pZalahYLe4o2iArcQHz3_VRsYd7frWN2kkF4ARwQl3UlOHH6jOSzT5h6PtnOJ1pDpIGME5NqG3TDvQnom5TAKW-XiZckk5lJAp3I51qGvDjve1AyZCRPfDHMsKAA`

const invalid_JWT_invalid_signature = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gU21pdGgiLCJhZG1pbiI6dHJ1ZSwiaWF0IjoxNTE2MjM5MDIyfQ.POstGetfAytaZS82wHcjoTyoqhMyxXiWdR7Nn7A29DNSl0EiXLdwJ6xC6AfgZWF1bOsS_TuYI3OG85AmiExREkrS6tDfTQ2B3WXlrr-wp5AokiRbz3_oB4OxG-W9KcEEbDRcZc0nH3L7LzYptiy1PtAylQGxHTWZXtGz4ht0bAecBgmpdgXMguEIcoqPJ1n3pIWk_dUZegpqx0Lka21H6XxUTxiy8OcaarA8zdnPUnV6AmNP3ecFawIFYdvJB_cm-GvpCSbr8G8y_Mllj8f4x9nBH8pQux89_6gUY618iYv7tuPWBFfEbLxtF2pZS6YC1aSfLQxeNe8djT9YjpvRZA`

const attack_JWT_alg_none = `eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.`
const attack_JWT_alg_symmetric = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.gqwq_XP_wdG5dqhJFFfUh4ico0HYmZ1-TgaMH3suC_I`

func TestVerifyJWT_no_JWT(t *testing.T) {
	r := httptest.NewRequest("GET", "https://doesntmatter.com/", nil)
	w := httptest.NewRecorder()

	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	// returns a function that creates an http.Handler instance
	middlewareFactory := JwtMiddleware(publicKey)

	// returns a handler instance (http.Handler)
	middleware := middlewareFactory(http.HandlerFunc(next))

	middleware.ServeHTTP(w, r)

	resp := w.Result()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	name, err := GetName(r.Context())
	require.Equal(t, "", name)
	require.EqualError(t, err, "failed to get name: failed to get token claims: no token in context")
	require.EqualError(t, errors.Unwrap(err), "no token in context")
	isAdmin, err := HasRole(r.Context(), "staff")
	require.False(t, isAdmin)
	require.EqualError(t, err, "failed to check for 'staff' role: failed to get token claims: no token in context")
	require.EqualError(t, errors.Unwrap(err), "no token in context")
}

func TestVerifyJWT_valid_JWT_staff(t *testing.T) {
	r := httptest.NewRequest("GET", "https://doesntmatter.com/", nil)
	r.Header.Set("authorization", "Bearer "+valid_JWT_is_staff)
	w := httptest.NewRecorder()

	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	// returns a function that creates an http.Handler instance
	middlewareFactory := JwtMiddleware(publicKey)

	// returns a handler instance (http.Handler)
	middleware := middlewareFactory(http.HandlerFunc(next))

	middleware.ServeHTTP(w, r)

	resp := w.Result()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	name, err := GetName(r.Context())
	require.Equal(t, "John Doe", name)
	require.Nil(t, err)
	isAdmin, err := HasRole(r.Context(), "staff")
	require.True(t, isAdmin)
	require.Nil(t, err)
}

func TestVerifyJWT_valid_JWT_not_staff(t *testing.T) {
	r := httptest.NewRequest("GET", "https://doesntmatter.com/", nil)
	r.Header.Set("authorization", "Bearer "+valid_JWT_is_not_staff)
	w := httptest.NewRecorder()

	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	// returns a function that creates an http.Handler instance
	middlewareFactory := JwtMiddleware(publicKey)

	// returns a handler instance (http.Handler)
	middleware := middlewareFactory(http.HandlerFunc(next))

	middleware.ServeHTTP(w, r)

	resp := w.Result()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	name, err := GetName(r.Context())
	require.Equal(t, "John Doe", name)
	require.Nil(t, err)
	isAdmin, err := HasRole(r.Context(), "staff")
	require.False(t, isAdmin)
	require.Nil(t, err)
}

func TestVerifyJWT_attack_none(t *testing.T) {
	r := httptest.NewRequest("GET", "https://doesntmatter.com/", nil)
	r.Header.Set("authorization", "Bearer "+attack_JWT_alg_none)
	w := httptest.NewRecorder()

	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	// returns a function that creates an http.Handler instance
	middlewareFactory := JwtMiddleware(publicKey)

	// returns a handler instance (http.Handler)
	middleware := middlewareFactory(http.HandlerFunc(next))

	middleware.ServeHTTP(w, r)

	resp := w.Result()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestVerifyJWT_attack_symmetric(t *testing.T) {
	r := httptest.NewRequest("GET", "https://doesntmatter.com/", nil)
	r.Header.Set("authorization", "Bearer "+attack_JWT_alg_symmetric)
	w := httptest.NewRecorder()

	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	// returns a function that creates an http.Handler instance
	middlewareFactory := JwtMiddleware(publicKey)

	// returns a handler instance (http.Handler)
	middleware := middlewareFactory(http.HandlerFunc(next))

	middleware.ServeHTTP(w, r)

	resp := w.Result()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestVerifyJWT_invalid_signature(t *testing.T) {
	r := httptest.NewRequest("GET", "https://doesntmatter.com/", nil)
	r.Header.Set("authorization", "Bearer "+invalid_JWT_invalid_signature)
	w := httptest.NewRecorder()

	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	// returns a function that creates an http.Handler instance
	middlewareFactory := JwtMiddleware(publicKey)

	// returns a handler instance (http.Handler)
	middleware := middlewareFactory(http.HandlerFunc(next))

	middleware.ServeHTTP(w, r)

	resp := w.Result()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
