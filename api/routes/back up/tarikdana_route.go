package routes

import (
	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/lib"
)

type TarikDanaRoutes struct {
	logger              lib.Logger
	handler             lib.HttpHandler
	tarikdanaController controllers.TarikDanaController
}

// NewTarikDanaRoutes creates new tarikdana routes
func NewTarikDanaRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	tarikdanaController controllers.TarikDanaController,
) TarikDanaRoutes {
	return TarikDanaRoutes{
		handler:             handler,
		logger:              logger,
		tarikdanaController: tarikdanaController,
	}
}

// Setup tarikdana routes
func (a TarikDanaRoutes) Setup() {
	a.logger.Zap.Info("Setting up tarikdana routes")
	api := a.handler.RouterV1.Group("/tarikdanas")
	{
		api.GET("", a.tarikdanaController.Query)
		api.GET(".all", a.tarikdanaController.GetAll)

		api.POST("", a.tarikdanaController.Create)
		api.GET("/:id", a.tarikdanaController.Get)
		api.PUT("/:id", a.tarikdanaController.Update)
		api.DELETE("/:id", a.tarikdanaController.Delete)
		api.PATCH("/:id/enable", a.tarikdanaController.Enable)
		api.PATCH("/:id/disable", a.tarikdanaController.Disable)
		api.GET("/upload/:id", a.tarikdanaController.GetFile)
		api.POST("/upload", a.tarikdanaController.UploadFile)
		api.DELETE("/upload/:id", a.tarikdanaController.RemoveFile)
	}
}
