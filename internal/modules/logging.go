package modules

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"openwt.com/go-arrp/internal/config"
)

func ProvideLogger(c *config.AppConfig) *zap.Logger {
	return c.Logger
}

var LoggingModule = fx.Options(
	fx.Provide(ProvideLogger),
)
