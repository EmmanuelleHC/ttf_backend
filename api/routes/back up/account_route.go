package routes

import (
	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/lib"
)

type AccountRoutes struct {
	logger            lib.Logger
	handler           lib.HttpHandler
	accountController controllers.AccountController
}

// NewAccountRoutes creates new account routes
func NewAccountRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	accountController controllers.AccountController,
) AccountRoutes {
	return AccountRoutes{
		handler:           handler,
		logger:            logger,
		accountController: accountController,
	}
}

// Setup account routes
func (a AccountRoutes) Setup() {
	a.logger.Zap.Info("Setting up account routes")
	api := a.handler.RouterV1.Group("/accounts")
	{
		api.GET("", a.accountController.Query)
		api.GET(".all", a.accountController.GetAll)

		api.POST("", a.accountController.Create)
		api.GET("/:id", a.accountController.Get)
		api.PUT("/:id", a.accountController.Update)
		api.DELETE("/:id", a.accountController.Delete)
		api.PATCH("/:id/enable", a.accountController.Enable)
		api.PATCH("/:id/disable", a.accountController.Disable)
	}
}
