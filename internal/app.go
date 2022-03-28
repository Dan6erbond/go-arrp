package internal

import (
	"go.uber.org/fx"
	"openwt.com/go-arrp/internal/modules"
)

func NewFxApp() *fx.App {
	app := fx.New(
		fx.Provide(NewGin),
		modules.AppModule,
		modules.JobsModule,
	)

	return app
}
