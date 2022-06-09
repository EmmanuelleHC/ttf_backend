package services

import "go.uber.org/fx"

// Module exports services present
var Module = fx.Options(
	fx.Provide(NewUserService),
	fx.Provide(NewRoleService),
	fx.Provide(NewMenuService),
	fx.Provide(NewCasbinService),
	fx.Provide(NewAuthService),
	fx.Provide(NewCompanyService),
	fx.Provide(NewBranchService),
	fx.Provide(NewDepartmentService),
	fx.Provide(NewAccountService),
	fx.Provide(NewCostCentreService),
	fx.Provide(NewSaldoService),
	fx.Provide(NewTrxService),
	fx.Provide(NewReportService),
	fx.Provide(NewKasbonService),
	fx.Provide(NewCounterService),
	fx.Provide(NewBKKDetailService),
	fx.Provide(NewBKKHeaderService),
	fx.Provide(NewEmployeeService),
	fx.Provide(NewTarikDanaService),
	fx.Provide(NewInvoiceHeaderService),
	fx.Provide(NewInvoiceDetailService),
)
