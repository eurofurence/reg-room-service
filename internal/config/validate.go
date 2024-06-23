package config

import (
	"errors"
	aulogging "github.com/StephanHCB/go-autumn-logging"
)

func (c *Config) Validate() error {
	ok := true

	if c.Server.Port <= 1024 || c.Server.Port > 65535 {
		aulogging.Logger.NoCtx().Warn().Print("server.port out of range")
		ok = false
	}

	if c.Service.MaxGroupSize < 1 {
		aulogging.Logger.NoCtx().Warn().Printf("need to set service.max_group_size")
		ok = false
	}

	// TODO more validation

	if ok {
		return nil
	} else {
		return errors.New("configuration validation error, see log output for details")
	}
}
