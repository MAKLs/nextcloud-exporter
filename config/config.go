package config

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	envPrefix  = "NC"
	configName = "config"
	configType = "yaml"
)

var (
	exporterConfig *Config
	configPaths    = []string{"."}
	defaults       = map[string]interface{}{
		"port":            9205,
		"token":           "",
		"url":             "http://localhost/",
		"exclude_php":     false,
		"exclude_strings": false,
		"filter":          []string{},
	}
)

type Config struct {
	Port           uint     `mapstructure:"port"`
	Url            url.URL  `mapstructure:"url"`
	Token          string   `mapstructure:"token"`
	FilterMetrics  []string `mapstructure:"filter"`
	ExcludePHP     bool     `mapstructure:"exclude_php"`
	ExcludeStrings bool     `mapstructure:"exclude_strings"`
}

func Notify(ch chan<- fsnotify.Event) {
	viper.OnConfigChange(func(in fsnotify.Event) {
		mustUnmarshalConfig()
		ch <- in
	})
}

func GetConfig() *Config {
	return exporterConfig
}

func decoderConfig(config *mapstructure.DecoderConfig) {
	config.ZeroFields = true
	config.ErrorUnused = true
}

func urlFromStringHookFunc() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t != reflect.TypeOf(url.URL{}) {
			return data, nil
		}

		return url.Parse(data.(string))
	}
}

func mustUnmarshalConfig() {
	if err := viper.Unmarshal(&exporterConfig, decoderConfig, viper.DecodeHook(urlFromStringHookFunc())); err != nil {
		panic(fmt.Sprintf("failed to parse config: %v", err))
	}
}

func init() {
	viper.SetEnvPrefix(envPrefix)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	for opt, val := range defaults {
		viper.SetDefault(opt, val)
	}

	viper.WatchConfig()
	viper.ReadInConfig()
	viper.AutomaticEnv()

	mustUnmarshalConfig()
}
