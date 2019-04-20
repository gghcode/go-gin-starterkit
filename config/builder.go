package config

import (
	"github.com/spf13/viper"
)

// Builder build configuration with some options.
type Builder struct {
	pipelines [](func(*viper.Viper) error)
}

// NewBuilder return new builder instance.
func NewBuilder() *Builder {
	return &Builder{}
}

// AddConfigFile read config file from filepath.
func (builder *Builder) AddConfigFile(filepath string, optional bool) *Builder {
	builder.pipelines = append(builder.pipelines, func(v *viper.Viper) error {

		return nil
	})

	return builder
}

// BindEnvs bind environment variables.
func (builder *Builder) BindEnvs(prefix string) *Builder {
	builder.pipelines = append(builder.pipelines, func(v *viper.Viper) error {
		return nil
	})

	return builder
}

// Build return new configuration instance.
func (builder *Builder) Build() (Configuration, error) {
	viperInstance := viper.New()

	for _, pipeline := range builder.pipelines {
		if err := pipeline(viperInstance); err != nil {
			return Configuration{}, err
		}
	}

	var result Configuration

	if err := viperInstance.Unmarshal(&result); err != nil {
		return Configuration{}, err
	}

	return result, nil
}
