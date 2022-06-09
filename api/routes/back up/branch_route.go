package routes

import (
	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/lib"
)

type BranchRoutes struct {
	logger           lib.Logger
	handler          lib.HttpHandler
	branchController controllers.BranchController
}

// NewBranchRoutes creates new branch routes
func NewBranchRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	branchController controllers.BranchController,
) BranchRoutes {
	return BranchRoutes{
		handler:          handler,
		logger:           logger,
		branchController: branchController,
	}
}

// Setup branch routes
func (a BranchRoutes) Setup() {
	a.logger.Zap.Info("Setting up branch routes")
	api := a.handler.RouterV1.Group("/branchs")
	{
		api.GET("", a.branchController.Query)
		api.GET(".all", a.branchController.GetAll)

		api.POST("", a.branchController.Create)
		api.GET("/:id", a.branchController.Get)
		api.PUT("/:id", a.branchController.Update)
		api.DELETE("/:id", a.branchController.Delete)
		api.PATCH("/:id/enable", a.branchController.Enable)
		api.PATCH("/:id/disable", a.branchController.Disable)
	}
}
