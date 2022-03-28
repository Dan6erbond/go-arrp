package modules

import (
	"go.uber.org/fx"
	"openwt.com/go-arrp/internal/controllers"
	"openwt.com/go-arrp/internal/services"
)

var JobsModule = fx.Options(
	fx.Provide(services.NewJobsService),
	fx.Provide(controllers.NewJobsController),
)
