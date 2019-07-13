package config

// Configuration is config type.
type Configuration struct {
	Addr     string         `mapstructure:"addr"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Jwt      JwtConfig      `mapstructure:"jwt"`
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

// JwtConfig is jwt config
type JwtConfig struct {
	SecretKey           string `mapstructure:"secret_key"`
	AccessExpiresInSec  int64  `mapstructure:"access_expires_sec"`
	RefreshExpiresInSec int64  `mapstructure:"refresh_expires_sec"`
}
