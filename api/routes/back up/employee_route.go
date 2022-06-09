package routes

import (
	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/lib"
)

type EmployeeRoutes struct {
	logger             lib.Logger
	handler            lib.HttpHandler
	employeeController controllers.EmployeeController
}

// NewEmployeeRoutes creates new employee routes
func NewEmployeeRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	employeeController controllers.EmployeeController,
) EmployeeRoutes {
	return EmployeeRoutes{
		handler:            handler,
		logger:             logger,
		employeeController: employeeController,
	}
}

// Setup employee routes
func (a EmployeeRoutes) Setup() {
	a.logger.Zap.Info("Setting up employee routes")
	api := a.handler.RouterV1.Group("/employees")
	{
		api.GET("", a.employeeController.Query)
		api.GET(".all", a.employeeController.GetAll)

		api.POST("", a.employeeController.Create)
		api.GET("/:id", a.employeeController.Get)
		api.PUT("/:id", a.employeeController.Update)
		api.DELETE("/:id", a.employeeController.Delete)
		api.PATCH("/:id/enable", a.employeeController.Enable)
		api.PATCH("/:id/disable", a.employeeController.Disable)
	}
}
