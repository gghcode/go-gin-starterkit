package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
)

type viperBuilder struct {
	viper *viper.Viper
}

func NewViperBuilder() viperBuilder {
	return viperBuilder{
		viper: viper.New(),
	}
}

func (builder *viperBuilder) BasePath(path string) Builder {
	builder.viper.AddConfigPath(path)
	return builder
}

func (builder *viperBuilder) KeyValue(key string, value interface{}) Builder {
	builder.viper.Set(key, value)

	return builder
}

func (builder *viperBuilder) JsonFile(path string) Builder {
	builder.viper.SetConfigType("json")
	builder.viper.SetConfigName(path)

	if err := builder.viper.MergeInConfig(); err != nil {
		log.Fatal(err)
	}

	return builder
}

func (builder *viperBuilder) EnvironmentVariables() Builder {
	panic(errors.New("Not implement method"))
}

func (builder *viperBuilder) Build() (Configuration, error) {
	result := Configuration{}

	if err := builder.viper.Unmarshal(&result); err != nil {
		return Configuration{}, err
	}

	return result, nil
}
