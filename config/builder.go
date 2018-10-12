package config

type Builder interface {
	Build() (Configuration, error)
}
