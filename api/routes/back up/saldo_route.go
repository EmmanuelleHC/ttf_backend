package routes

import (
	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/lib"
)

type SaldoRoutes struct {
	logger          lib.Logger
	handler         lib.HttpHandler
	saldoController controllers.SaldoController
}

// NewSaldoRoutes creates new saldo routes
func NewSaldoRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	saldoController controllers.SaldoController,
) SaldoRoutes {
	return SaldoRoutes{
		handler:         handler,
		logger:          logger,
		saldoController: saldoController,
	}
}

// Setup saldo routes
func (a SaldoRoutes) Setup() {
	a.logger.Zap.Info("Setting up saldo routes")
	api := a.handler.RouterV1.Group("/saldos")
	{
		api.GET("", a.saldoController.Query)
		api.GET(".all", a.saldoController.GetAll)

		api.POST("", a.saldoController.Create)
		api.GET("/:id", a.saldoController.Get)
		api.GET("/user", a.saldoController.GetByUser)
		api.PUT("/:id", a.saldoController.Update)
		api.DELETE("/:id", a.saldoController.Delete)
		api.PATCH("/:id/enable", a.saldoController.Enable)
		api.PATCH("/:id/disable", a.saldoController.Disable)
	}
}
