package config

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AppConfig struct {
	Env     string
	Cadence CadenceConfig
	Logger  *zap.Logger
}

func (h *AppConfig) Setup(c ...*AppConfig) {
	viper.SetConfigName("application")
	viper.AddConfigPath("configs")
	viper.SetConfigType("yml")

	viper.AutomaticEnv() // This allows viper to read variables from the environment variables if they exists.

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&h)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	if len(c) > 0 {
		copier.Copy(&h, c[0])
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	h.Logger = logger

	logger.Debug("Finished loading Configuration!")
}

func ProvideConfig(c ...*AppConfig) *AppConfig {
	var config AppConfig
	config.Setup(c...)
	return &config
}

var Module = fx.Options(
	fx.Provide(ProvideConfig),
)
