package modules

import (
	"go.uber.org/fx"
	"openwt.com/go-arrp/internal/cadence"
	"openwt.com/go-arrp/internal/config"
)

func ProvideCadenceAdapter(c *config.CadenceConfig) *cadence.CadenceAdapter {
	var cadenceClient cadence.CadenceAdapter
	cadenceClient.Setup(c)
	return &cadenceClient
}

var CadenceModule = fx.Options(
	fx.Provide(ProvideCadenceAdapter),
)
