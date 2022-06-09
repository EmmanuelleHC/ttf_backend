package routes

import (
	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/lib"
)

type BKKDetailRoutes struct {
	logger              lib.Logger
	handler             lib.HttpHandler
	bkkdetailController controllers.BKKDetailController
}

// NewBKKDetailRoutes creates new bkkdetail routes
func NewBKKDetailRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	bkkdetailController controllers.BKKDetailController,
) BKKDetailRoutes {
	return BKKDetailRoutes{
		handler:             handler,
		logger:              logger,
		bkkdetailController: bkkdetailController,
	}
}

// Setup bkkdetail routes
func (a BKKDetailRoutes) Setup() {
	a.logger.Zap.Info("Setting up bkkdetail routes")
	api := a.handler.RouterV1.Group("/bkkdetails")
	{
		api.GET("", a.bkkdetailController.Query)
		api.GET(".all", a.bkkdetailController.GetAll)

		api.POST("", a.bkkdetailController.Create)
		api.GET("/:id", a.bkkdetailController.Get)
		api.PUT("/:id", a.bkkdetailController.Update)
		api.DELETE("/:id", a.bkkdetailController.Delete)
		api.PATCH("/:id/enable", a.bkkdetailController.Enable)
		api.PATCH("/:id/disable", a.bkkdetailController.Disable)
		api.GET("/upload/:id", a.bkkdetailController.GetFile)
		api.POST("/upload", a.bkkdetailController.UploadFile)
		api.DELETE("/upload/:id", a.bkkdetailController.RemoveFile)
	}
}
