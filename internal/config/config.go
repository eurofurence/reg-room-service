package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var appConfig *Config

type (
	DatabaseType string
	LogStyle     string
)

const (
	Mysql    DatabaseType = "mysql"
	Inmemory DatabaseType = "inmemory"

	Plain LogStyle = "plain"
	ECS   LogStyle = "ecs" // default
)

type (
	// Config is the root configuration type
	// that holds all other subconfiguration types.
	Config struct {
		Service  ServiceConfig  `yaml:"service"`
		Server   ServerConfig   `yaml:"server"`
		Database DatabaseConfig `yaml:"database"`
		Security SecurityConfig `yaml:"security"`
		Logging  LoggingConfig  `yaml:"logging"`
		GoLive   GoLiveConfig   `yaml:"go_live"`
	}

	// ServiceConfig contains configuration values
	// for service related tasks. E.g. URL to attendee service.
	ServiceConfig struct {
		AttendeeServiceURL string `yaml:"attendee_service_url"`
	}

	// ServerConfig contains all values for
	// http releated configuration.
	ServerConfig struct {
		BaseAddress  string `yaml:"address"`
		Port         int    `yaml:"port"`
		ReadTimeout  int    `yaml:"read_timeout_seconds"`
		WriteTimeout int    `yaml:"write_timeout_seconds"`
		IdleTimeout  int    `yaml:"idle_timeout_seconds"`
	}

	// DatabaseConfig configures which db to use (mysql, inmemory)
	// and how to connect to it (needed for mysql only).
	DatabaseConfig struct {
		Use        DatabaseType `yaml:"use"`
		Username   string       `yaml:"username"`
		Password   string       `yaml:"password"`
		Database   string       `yaml:"database"`
		Parameters []string     `yaml:"parameters"`
	}

	// SecurityConfig configures everything related to security.
	SecurityConfig struct {
		Fixed        FixedTokenConfig    `yaml:"fixed_token"`
		Oidc         OpenIDConnectConfig `yaml:"oidc"`
		Cors         CorsConfig          `yaml:"cors"`
		RequireLogin bool                `yaml:"require_login_for_reg"`
	}

	FixedTokenConfig struct {
		API string `yaml:"api"` // shared-secret for server-to-server backend authentication
	}

	OpenIDConnectConfig struct {
		IDTokenCookieName     string   `yaml:"id_token_cookie_name"`     // optional, but must both be set, then tokens are read from cookies
		AccessTokenCookieName string   `yaml:"access_token_cookie_name"` // optional, but must both be set, then tokens are read from cookies
		TokenPublicKeysPEM    []string `yaml:"token_public_keys_PEM"`    // a list of public RSA keys in PEM format, see https://github.com/Jumpy-Squirrel/jwks2pem for obtaining PEM from openid keyset endpoint
		AdminGroup            string   `yaml:"admin_group"`              // the group claim that supplies admin rights
		AuthService           string   `yaml:"auth_service"`             // base url, usually http://localhost:nnnn, will skip userinfo checks if unset
		Audience              string   `yaml:"audience"`
		Issuer                string   `yaml:"issuer"`
	}

	CorsConfig struct {
		DisableCors bool   `yaml:"disable"`
		AllowOrigin string `yaml:"allow_origin"`
	}

	// LoggingConfig configures logging.
	LoggingConfig struct {
		Style    LogStyle `yaml:"style"`
		Severity string   `yaml:"severity"`
	}

	GoLiveConfig struct {
		Public GoLiveConfigPerGroup `yaml:"public"`
		Staff  GoLiveConfigPerGroup `yaml:"staff"`
	}

	GoLiveConfigPerGroup struct {
		StartISODatetime string `yaml:"start_iso_datetime"`
		BookingCode      string `yaml:"booking_code"`
		Group            string `yaml:"group"`
	}
)

// UnmarshalFromYamlConfiguration decodes yaml data from an `io.Reader` interface.
func UnmarshalFromYamlConfiguration(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, errors.New("no config path provided")
	}

	f, err := os.Open(filepath.Clean(configPath))
	if err != nil {
		return nil, err
	}

	d := yaml.NewDecoder(f)
	d.KnownFields(true) // strict

	var conf Config

	if err := d.Decode(&conf); err != nil {
		return nil, err
	}

	appConfig = &conf
	return &conf, nil
}

func GetApplicationConfig() (*Config, error) {
	if appConfig == nil {
		return nil, errors.New("config was not yet loaded")
	}

	return appConfig, nil
}
