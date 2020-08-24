package configs

import (
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/whale-team/whaleEcho/internal/pkg/repository/db"
	"github.com/whale-team/whaleEcho/pkg/stanclient"
	"go.uber.org/fx"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/rs/zerolog/log"
	"github.com/vicxu416/goinfra/zlogging"
)

var _config Configuration

// Configuration represent app configuration
type Configuration struct {
	fx.Out
	WsServer struct {
		Addr string `yaml:"addr"`
		Port string `yaml:"port"`
	} `yaml:"ws_server"`

	Log   zlogging.Config   `yaml:"log"`
	Stan  stanclient.Config `yaml:"stan"`
	Redis db.RedisConfig    `yaml:"redis"`
}

// Empty check if configuration is empty
func (c Configuration) Empty() bool {
	return reflect.DeepEqual(c, Configuration{})
}

// InitConfiguration initialize configuration
func InitConfiguration() (Configuration, error) {
	if !_config.Empty() {
		return _config, nil
	}

	viper.AutomaticEnv()

	configPath := viper.GetString("CONFIGPATH")
	if configPath == "" {
		_, f, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(f)
		configPath = filepath.Join(basepath, "/")
	}

	configName := viper.GetString("CONFIG_NAME")
	if configName == "" {
		configName = "app"
	}

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")

	var config Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Msgf("Error reading config file, %s", err)
		return config, err
	}

	err := viper.Unmarshal(&config, func(c *mapstructure.DecoderConfig) {
		c.TagName = "yaml"
	})
	if err != nil {
		log.Error().Msgf("unable to decode into struct, %v", err)
		return config, err
	}

	_config = config
	return _config, nil
}
