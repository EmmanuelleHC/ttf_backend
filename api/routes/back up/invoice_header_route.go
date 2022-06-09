package routes

import (
	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/lib"
)

type InvoiceHeaderRoutes struct {
	logger                  lib.Logger
	handler                 lib.HttpHandler
	invoiceheaderController controllers.InvoiceHeaderController
}

// NewInvoiceHeaderRoutes creates new invoiceheader routes
func NewInvoiceHeaderRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	invoiceheaderController controllers.InvoiceHeaderController,
) InvoiceHeaderRoutes {
	return InvoiceHeaderRoutes{
		handler:                 handler,
		logger:                  logger,
		invoiceheaderController: invoiceheaderController,
	}
}

// Setup invoiceheader routes
func (a InvoiceHeaderRoutes) Setup() {
	a.logger.Zap.Info("Setting up invoiceheader routes")
	api := a.handler.RouterV1.Group("/invoiceheaders")
	{
		api.GET("", a.invoiceheaderController.Query)
		api.GET(".all", a.invoiceheaderController.GetAll)

		api.POST("", a.invoiceheaderController.Create)
		api.GET("/:id", a.invoiceheaderController.Get)
		api.PUT("/:id", a.invoiceheaderController.Update)
		api.DELETE("/:id", a.invoiceheaderController.Delete)
		api.PATCH("/:id/enable", a.invoiceheaderController.Enable)
		api.PATCH("/:id/disable", a.invoiceheaderController.Disable)

		api.POST("/approve", a.invoiceheaderController.StatusApprove)
		api.POST("/reject", a.invoiceheaderController.StatusReject)
		api.POST("/approvefinal", a.invoiceheaderController.StatusApproveFinal)
		api.POST("/rejectfinal", a.invoiceheaderController.StatusRejectFinal)
	}
}
