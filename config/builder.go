package config

type Builder interface {
	SetBasePath(path string) Builder
	SetValue(key string, value interface{}) Builder

	AddJsonFile(path string) Builder
	AddEnvironmentVariables() Builder

	Build() (Configuration, error)
}
