package config

import (
	"errors"
	aulogging "github.com/StephanHCB/go-autumn-logging"
)

func (c *Config) Validate() error {
	ok := true

	if c.Server.Port <= 1024 || c.Server.Port > 65535 {
		aulogging.Logger.NoCtx().Error().Print("server.port out of range")
		ok = false
	}

	// TODO more validation

	if ok {
		return nil
	} else {
		return errors.New("configuration validation error, see log output for details")
	}
}
