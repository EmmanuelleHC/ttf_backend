package routes

import (
	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/lib"
)

type BKKHeaderRoutes struct {
	logger              lib.Logger
	handler             lib.HttpHandler
	bkkheaderController controllers.BKKHeaderController
}

// NewBKKHeaderRoutes creates new bkkheader routes
func NewBKKHeaderRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	bkkheaderController controllers.BKKHeaderController,
) BKKHeaderRoutes {
	return BKKHeaderRoutes{
		handler:             handler,
		logger:              logger,
		bkkheaderController: bkkheaderController,
	}
}

// Setup bkkheader routes
func (a BKKHeaderRoutes) Setup() {
	a.logger.Zap.Info("Setting up bkkheader routes")
	api := a.handler.RouterV1.Group("/bkkheaders")
	{
		api.GET("", a.bkkheaderController.Query)
		api.GET(".all", a.bkkheaderController.GetAll)

		api.POST("", a.bkkheaderController.Create)
		api.GET("/:id", a.bkkheaderController.Get)
		api.PUT("/:id", a.bkkheaderController.Update)
		api.DELETE("/:id", a.bkkheaderController.Delete)
		api.PATCH("/:id/statuspaid", a.bkkheaderController.StatusPaid)
		api.POST("/approve", a.bkkheaderController.StatusApprove)
		api.POST("/reject", a.bkkheaderController.StatusReject)
	}
}
