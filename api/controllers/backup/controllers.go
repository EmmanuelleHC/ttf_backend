package controllers

import "go.uber.org/fx"

// Module exported for initializing application
var Module = fx.Options(
	fx.Provide(NewPublicController),
	fx.Provide(NewCaptchaController),
	fx.Provide(NewUserController),
	fx.Provide(NewRoleController),
	fx.Provide(NewMenuController),
	fx.Provide(NewCompanyController),
	fx.Provide(NewBranchController),
	fx.Provide(NewDepartmentController),
	fx.Provide(NewAccountController),
	fx.Provide(NewCostCentreController),
	fx.Provide(NewSaldoController),
	fx.Provide(NewTrxController),
	fx.Provide(NewKasbonController),
	fx.Provide(NewBKKDetailController),
	fx.Provide(NewBKKHeaderController),
	fx.Provide(NewEmployeeController),
	fx.Provide(NewTarikDanaController),
	fx.Provide(NewInvoiceHeaderController),
	fx.Provide(NewInvoiceDetailController),
)
