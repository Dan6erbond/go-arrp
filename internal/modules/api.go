package modules

import (
	"go.uber.org/fx"
	"openwt.com/go-arrp/internal/router"
)

var ApiModule = fx.Options(
	fx.Invoke(router.RegisterRoutes),
)
