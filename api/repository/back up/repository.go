package repository

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Options(
	fx.Provide(NewUserRepository),
	fx.Provide(NewUserRoleRepository),
	fx.Provide(NewRoleRepository),
	fx.Provide(NewRoleMenuRepository),
	fx.Provide(NewMenuRepository),
	fx.Provide(NewMenuActionRepository),
	fx.Provide(NewMenuActionResourceRepository),
	fx.Provide(NewCompanyRepository),
	fx.Provide(NewBranchRepository),
	fx.Provide(NewDepartmentRepository),
	fx.Provide(NewAccountRepository),
	fx.Provide(NewCostCentreRepository),
	fx.Provide(NewSaldoRepository),
	fx.Provide(NewSaldoHistoryRepository),
	fx.Provide(NewSaldoMonthRepository),
	fx.Provide(NewTrxRepository),
	fx.Provide(NewKasbonRepository),
	fx.Provide(NewCounterRepository),
	fx.Provide(NewBKKDetailRepository),
	fx.Provide(NewBKKHeaderRepository),
	fx.Provide(NewEmployeeRepository),
	fx.Provide(NewTarikDanaRepository),
	fx.Provide(NewInvoiceHeaderRepository),
	fx.Provide(NewInvoiceDetailRepository),
)
