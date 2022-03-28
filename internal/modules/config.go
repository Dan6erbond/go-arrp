package modules

import (
	"go.uber.org/fx"
	"openwt.com/go-arrp/internal/config"
)

func ProvideAppConfig() *config.AppConfig {
	appConfig := config.ProvideConfig()
	return appConfig
}

func ProvideCadenceConfig(c *config.AppConfig) *config.CadenceConfig {
	return &c.Cadence
}

var ConfigModule = fx.Options(
	fx.Provide(ProvideAppConfig),
	fx.Provide(ProvideCadenceConfig),
)
