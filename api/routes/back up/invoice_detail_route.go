package routes

import (
	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/lib"
)

type InvoiceDetailRoutes struct {
	logger                  lib.Logger
	handler                 lib.HttpHandler
	invoicedetailController controllers.InvoiceDetailController
}

// NewInvoiceDetailRoutes creates new invoicedetail routes
func NewInvoiceDetailRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	invoicedetailController controllers.InvoiceDetailController,
) InvoiceDetailRoutes {
	return InvoiceDetailRoutes{
		handler:                 handler,
		logger:                  logger,
		invoicedetailController: invoicedetailController,
	}
}

// Setup invoicedetail routes
func (a InvoiceDetailRoutes) Setup() {
	a.logger.Zap.Info("Setting up invoicedetail routes")
	api := a.handler.RouterV1.Group("/invoicedetails")
	{
		api.GET("", a.invoicedetailController.Query)
		api.GET("/bkk", a.invoicedetailController.QueryBkk)
		api.GET(".all", a.invoicedetailController.GetAll)

		api.POST("", a.invoicedetailController.Create)
		api.GET("/:id", a.invoicedetailController.Get)
		api.PUT("/:id", a.invoicedetailController.Update)
		api.DELETE("/:id", a.invoicedetailController.Delete)
		api.PATCH("/:id/enable", a.invoicedetailController.Enable)
		api.PATCH("/:id/disable", a.invoicedetailController.Disable)
	}
}
