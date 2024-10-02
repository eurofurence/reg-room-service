package config

import "os"

const (
	envDbPassword = "REG_SECRET_DB_PASSWORD"
	envApiToken   = "REG_SECRET_API_TOKEN"
)

func (c *Config) ApplyEnvironmentOverrides() {
	if dbPassword := os.Getenv(envDbPassword); dbPassword != "" {
		c.Database.Password = dbPassword
	}
	if apiToken := os.Getenv(envApiToken); apiToken != "" {
		c.Security.Fixed.API = apiToken
	}
}
