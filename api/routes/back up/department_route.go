package routes

import (
	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/lib"
)

type DepartmentRoutes struct {
	logger               lib.Logger
	handler              lib.HttpHandler
	departmentController controllers.DepartmentController
}

// NewDepartmentRoutes creates new department routes
func NewDepartmentRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	departmentController controllers.DepartmentController,
) DepartmentRoutes {
	return DepartmentRoutes{
		handler:              handler,
		logger:               logger,
		departmentController: departmentController,
	}
}

// Setup department routes
func (a DepartmentRoutes) Setup() {
	a.logger.Zap.Info("Setting up department routes")
	api := a.handler.RouterV1.Group("/departments")
	{
		api.GET("", a.departmentController.Query)
		api.GET(".all", a.departmentController.GetAll)

		api.POST("", a.departmentController.Create)
		api.GET("/:id", a.departmentController.Get)
		api.PUT("/:id", a.departmentController.Update)
		api.DELETE("/:id", a.departmentController.Delete)
		api.PATCH("/:id/enable", a.departmentController.Enable)
		api.PATCH("/:id/disable", a.departmentController.Disable)
	}
}
