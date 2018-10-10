package config

type Configuration struct {
	ListenPort            int    `mapstructure:"LISTEN_PORT"`
	MongoDbName           string `mapstructure:"MONGO_DB_NAME"`
	MongoConnectionString string `mapstructure:"MONGO_CONN_STR"`
}
