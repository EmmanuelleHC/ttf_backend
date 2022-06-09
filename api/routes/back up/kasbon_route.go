package routes

import (
	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/lib"
)

type KasbonRoutes struct {
	logger           lib.Logger
	handler          lib.HttpHandler
	kasbonController controllers.KasbonController
}

// NewKasbonRoutes creates new kasbon routes
func NewKasbonRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	kasbonController controllers.KasbonController,
) KasbonRoutes {
	return KasbonRoutes{
		handler:          handler,
		logger:           logger,
		kasbonController: kasbonController,
	}
}

// Setup kasbon routes
func (a KasbonRoutes) Setup() {
	a.logger.Zap.Info("Setting up kasbon routes")
	api := a.handler.RouterV1.Group("/kasbons")
	{
		api.GET("", a.kasbonController.Query)
		api.GET(".all", a.kasbonController.GetAll)

		api.POST("", a.kasbonController.Create)
		api.GET("/:id", a.kasbonController.Get)
		api.PUT("/:id", a.kasbonController.Update)
		api.DELETE("/:id", a.kasbonController.Delete)
		api.GET("/upload/:id", a.kasbonController.GetFile)
		api.POST("/upload", a.kasbonController.UploadFile)
		api.DELETE("/upload/:id", a.kasbonController.RemoveFile)
	}
}
