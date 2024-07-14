package config

func (c *Config) AddDefaults() {
	if c.Server.Port == 0 {
		c.Server.Port = 8081
	}
	if c.Server.IdleTimeout <= 0 {
		c.Server.IdleTimeout = 30
	}
	if c.Server.ReadTimeout <= 0 {
		c.Server.ReadTimeout = 30
	}
	if c.Server.WriteTimeout <= 0 {
		c.Server.WriteTimeout = 30
	}
}
