package routes

import (
	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/lib"
)

type CompanyRoutes struct {
	logger            lib.Logger
	handler           lib.HttpHandler
	companyController controllers.CompanyController
}

// NewCompanyRoutes creates new company routes
func NewCompanyRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	companyController controllers.CompanyController,
) CompanyRoutes {
	return CompanyRoutes{
		handler:           handler,
		logger:            logger,
		companyController: companyController,
	}
}

// Setup company routes
func (a CompanyRoutes) Setup() {
	a.logger.Zap.Info("Setting up company routes")
	api := a.handler.RouterV1.Group("/companies")
	{
		api.GET("", a.companyController.Query)
		api.GET(".all", a.companyController.GetAll)

		api.POST("", a.companyController.Create)
		api.GET("/:id", a.companyController.Get)
		api.PUT("/:id", a.companyController.Update)
		api.DELETE("/:id", a.companyController.Delete)
		api.PATCH("/:id/enable", a.companyController.Enable)
		api.PATCH("/:id/disable", a.companyController.Disable)
	}
}
