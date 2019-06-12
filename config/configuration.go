package config

// Configuration is config type.
type Configuration struct {
	Addr     string         `mapstructure:"addr"`
	Postgres PostgresConfig `mapstructure:"postgres"`
}

// PostgresConfig is postgres config
type PostgresConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
}
