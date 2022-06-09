package routes

import (
	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/lib"
)

type TrxRoutes struct {
	logger        lib.Logger
	handler       lib.HttpHandler
	trxController controllers.TrxController
}

// NewTrxRoutes creates new trx routes
func NewTrxRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	trxController controllers.TrxController,
) TrxRoutes {
	return TrxRoutes{
		handler:       handler,
		logger:        logger,
		trxController: trxController,
	}
}

// Setup trx routes
func (a TrxRoutes) Setup() {
	a.logger.Zap.Info("Setting up trx routes")
	api := a.handler.RouterV1.Group("/trxs")
	{
		api.GET("", a.trxController.Query)
		api.GET(".all", a.trxController.GetAll)

		api.POST("", a.trxController.Create)
		api.GET("/:id", a.trxController.Get)
		api.PUT("/:id", a.trxController.Update)
		api.DELETE("/:id", a.trxController.Delete)
		api.PATCH("/:id/enable", a.trxController.Enable)
		api.PATCH("/:id/disable", a.trxController.Disable)
	}
}
