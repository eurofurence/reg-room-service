package rbac

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"

	"github.com/eurofurence/reg-room-service/internal/config"
	"github.com/eurofurence/reg-room-service/internal/web/common"
)

// note: there is a TestMain that loads configuration

func TestMain(m *testing.M) {
	_, err := config.UnmarshalFromYamlConfiguration(filepath.Join("..", "..", "..", "docs", "config.example.yaml"))
	if err != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())

}

func TestNewRBACValidator(t *testing.T) {
	type args struct {
		inputJWT               string
		inputAPIKey            string
		inputClaims            *common.AllClaims
		includeAdminHeader     bool
		customAdminHeaderValue string
	}

	type expected struct {
		subject          string
		roles            []string
		isAdmin          bool
		isAPITokenCall   bool
		isRegisteredUser bool
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "Should create manager with valid API token",
			args: args{
				inputJWT:    "",
				inputAPIKey: "api-token",
				inputClaims: nil,
			},
			expected: expected{
				isAPITokenCall: true,
			},
		},
		{
			name: "Should create manager with admin role",
			args: args{
				inputJWT:    "valid",
				inputAPIKey: "",
				inputClaims: &common.AllClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						Subject: "123456",
					},
					CustomClaims: common.CustomClaims{
						Groups: []string{"admin", "test"},
						Name:   "Peter",
						EMail:  "peter@peter.eu",
					},
				},
				includeAdminHeader: true,
			},
			expected: expected{
				isAdmin: true,
				subject: "123456",
				roles:   []string{"admin", "test"},
			},
		},
		// TODO remove test case after 2FA is available
		// See reference https://github.com/eurofurence/reg-payment-service/issues/57
		{
			name: "Should not create manager with admin role when no admin header is set",
			args: args{
				inputJWT:    "valid",
				inputAPIKey: "",
				inputClaims: &common.AllClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						Subject: "123456",
					},
					CustomClaims: common.CustomClaims{
						Groups: []string{"admin", "test"},
						Name:   "Peter",
						EMail:  "peter@peter.eu",
					},
				},
				includeAdminHeader: false,
			},
			expected: expected{
				isAdmin:          false,
				isRegisteredUser: true,
				subject:          "123456",
				roles:            []string{"admin", "test"},
			},
		},
		// TODO remove test case after 2FA is available
		// See reference https://github.com/eurofurence/reg-payment-service/issues/57
		{
			name: "Should not create manager with admin role when no valid admin header is set",
			args: args{
				inputJWT:    "valid",
				inputAPIKey: "",
				inputClaims: &common.AllClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						Subject: "123456",
					},
					CustomClaims: common.CustomClaims{
						Groups: []string{"admin", "test"},
						Name:   "Peter",
						EMail:  "peter@peter.eu",
					},
				},
				includeAdminHeader:     true,
				customAdminHeaderValue: "test-12345",
			},
			expected: expected{
				isAdmin:          false,
				isRegisteredUser: true,
				subject:          "123456",
				roles:            []string{"admin", "test"},
			},
		},
		{
			name: "Should create manager with registered user role",
			args: args{
				inputJWT:    "valid",
				inputAPIKey: "",
				inputClaims: &common.AllClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						Subject: "123456",
					},
					CustomClaims: common.CustomClaims{
						Groups: []string{"staff", "test"},
						Name:   "Peter",
						EMail:  "peter@peter.eu",
					},
				},
			},
			expected: expected{
				isRegisteredUser: true,
				subject:          "123456",
				roles:            []string{"staff", "test"},
			},
		},
		{
			name: "Should return empty manager if no tokens provided",
			args: args{
				inputJWT:    "",
				inputAPIKey: "",
				inputClaims: nil,
			},
			expected: expected{},
		},
		{
			name: "Should return empty manager if no tokens provided",
			args: args{
				inputJWT:    "",
				inputAPIKey: "",
				inputClaims: nil,
			},
			expected: expected{},
		},
		{
			name: "Should be invalid if registered user has no subject assigned",
			args: args{
				inputJWT:    "valid",
				inputAPIKey: "",
				inputClaims: &common.AllClaims{
					CustomClaims: common.CustomClaims{
						Groups: []string{""},
					},
				},
			},
			expected: expected{
				roles: []string{""},
			},
		},
		{
			name: "API key should dominate over JWT token",
			args: args{
				inputJWT:    "valid",
				inputAPIKey: "also valid",
				inputClaims: &common.AllClaims{
					CustomClaims: common.CustomClaims{
						Groups: []string{""},
					},
				},
			},
			expected: expected{
				isAPITokenCall: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.args.inputAPIKey != "" {
				ctx = context.WithValue(ctx, common.CtxKeyAPIKey{}, tt.args.inputAPIKey)
			}

			if tt.args.inputJWT != "" {
				ctx = context.WithValue(ctx, common.CtxKeyIDToken{}, tt.args.inputJWT)
			}

			if tt.args.inputClaims != nil {
				ctx = context.WithValue(ctx, common.CtxKeyClaims{}, tt.args.inputClaims)
				if tt.args.includeAdminHeader {
					ctx = context.WithValue(ctx, common.CtxKeyAdminHeader{}, coalesce(tt.args.customAdminHeaderValue, "available"))
				}
			}

			mgr, err := NewValidator(ctx)
			require.Nil(t, err)

			require.Equal(t, tt.expected.isAdmin, mgr.IsAdmin())
			require.Equal(t, tt.expected.isAPITokenCall, mgr.IsAPITokenCall())
			require.Equal(t, tt.expected.isRegisteredUser, mgr.IsRegisteredUser())
			require.Equal(t, tt.expected.roles, mgr.Groups())
			require.Equal(t, tt.expected.subject, mgr.Subject())

		})
	}
}

func coalesce(input, defaultValue string) string {
	if input == "" {
		return defaultValue
	}

	return input
}
