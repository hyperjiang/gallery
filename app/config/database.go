package config

// DatabaseSettings - settings for open a database
type DatabaseSettings struct {
	URI             string `yaml:"uri"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

// DatabaseConfig - configs for the databases
type DatabaseConfig struct {
	Main DatabaseSettings `yaml:"main"`
}
