package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AppConfig struct {
	Env            string
	WorkerTaskList string
	Cadence        CadenceConfig
	Logger         *zap.Logger
}

func (h *AppConfig) Setup() {
	viper.SetConfigName("application")
	viper.AddConfigPath("configs") // These two lines will make sure viper pulls the config from app/resources/application.yml
	viper.SetConfigType("yml")

	viper.AutomaticEnv() // This allows viper to read variables from the environment variables if they exists.

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&h)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	h.Logger = logger

	logger.Debug("Finished loading Configuration!")
}
