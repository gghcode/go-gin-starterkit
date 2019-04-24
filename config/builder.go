package config

import (
	"os"
	"reflect"
	"strings"

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

// AddConfigFile read config file from filePath.
func (builder *Builder) AddConfigFile(filePath string, optional bool) *Builder {
	builder.pipelines = append(builder.pipelines, func(viperObj *viper.Viper) error {
		viperObj.SetConfigFile(filePath)

		err := viperObj.MergeInConfig()
		if err != nil {
			if _, ok := err.(*os.PathError); ok && !optional {
				return err
			}
		}

		return nil
	})

	return builder
}

// BindEnvs bind environment variables.
func (builder *Builder) BindEnvs(prefix string) *Builder {
	builder.pipelines = append(builder.pipelines, func(viperObj *viper.Viper) error {
		viperObj.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viperObj.SetEnvPrefix(prefix)

		bindEnvsToViper(viperObj, Configuration{})

		return nil
	})

	return builder
}

func bindEnvsToViper(viperObj *viper.Viper, iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}

		switch v.Kind() {
		case reflect.Struct:
			bindEnvsToViper(viperObj, v.Interface(), append(parts, tv)...)
		default:
			viperObj.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}

// Build return new configuration instance.
func (builder *Builder) Build() (Configuration, error) {
	viperObj := viper.New()

	for _, pipeline := range builder.pipelines {
		if err := pipeline(viperObj); err != nil {
			return Configuration{}, err
		}
	}

	var result Configuration

	if err := viperObj.Unmarshal(&result); err != nil {
		return Configuration{}, err
	}

	return result, nil
}
