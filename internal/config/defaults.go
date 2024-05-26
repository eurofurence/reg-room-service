package config

func (c *Config) AddDefaults() {
	if c.Server.Port == 0 {
		c.Server.Port = 8081
	}
}
