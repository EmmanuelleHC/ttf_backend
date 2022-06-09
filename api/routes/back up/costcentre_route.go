package routes

import (
	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/lib"
)

type CostCentreRoutes struct {
	logger               lib.Logger
	handler              lib.HttpHandler
	costcentreController controllers.CostCentreController
}

// NewCostCentreRoutes creates new costcentre routes
func NewCostCentreRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	costcentreController controllers.CostCentreController,
) CostCentreRoutes {
	return CostCentreRoutes{
		handler:              handler,
		logger:               logger,
		costcentreController: costcentreController,
	}
}

// Setup costcentre routes
func (a CostCentreRoutes) Setup() {
	a.logger.Zap.Info("Setting up costcentre routes")
	api := a.handler.RouterV1.Group("/costcentres")
	{
		api.GET("", a.costcentreController.Query)
		api.GET(".all", a.costcentreController.GetAll)

		api.POST("", a.costcentreController.Create)
		api.GET("/:id", a.costcentreController.Get)
		api.PUT("/:id", a.costcentreController.Update)
		api.DELETE("/:id", a.costcentreController.Delete)
		api.PATCH("/:id/enable", a.costcentreController.Enable)
		api.PATCH("/:id/disable", a.costcentreController.Disable)
	}
}
