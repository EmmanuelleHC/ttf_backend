package routes

import "go.uber.org/fx"

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewPprofRoutes),
	fx.Provide(NewSwaggerRoutes),
	fx.Provide(NewPublicRoutes),
	fx.Provide(NewUserRoutes),
	fx.Provide(NewRoleRoutes),
	fx.Provide(NewMenuRoutes),
	fx.Provide(NewRoutes),
	fx.Provide(NewCompanyRoutes),
	fx.Provide(NewBranchRoutes),
	fx.Provide(NewDepartmentRoutes),
	fx.Provide(NewAccountRoutes),
	fx.Provide(NewCostCentreRoutes),
	fx.Provide(NewSaldoRoutes),
	fx.Provide(NewTrxRoutes),
	fx.Provide(NewKasbonRoutes),
	fx.Provide(NewBKKDetailRoutes),
	fx.Provide(NewBKKHeaderRoutes),
	fx.Provide(NewEmployeeRoutes),
	fx.Provide(NewTarikDanaRoutes),
	fx.Provide(NewInvoiceHeaderRoutes),
	fx.Provide(NewInvoiceDetailRoutes),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	pprofRoutes PprofRoutes,
	swaggerRoutes SwaggerRoutes,
	publicRoutes PublicRoutes,
	userRoutes UserRoutes,
	roleRoutes RoleRoutes,
	menuRoutes MenuRoutes,
	companyRoutes CompanyRoutes,
	branchRoutes BranchRoutes,
	departmentRoutes DepartmentRoutes,
	accountRoutes AccountRoutes,
	costcentreRoutes CostCentreRoutes,
	saldoRoutes SaldoRoutes,
	trxRoutes TrxRoutes,
	kasbonRoutes KasbonRoutes,
	bkkdetailRoutes BKKDetailRoutes,
	bkkHeaderRoutes BKKHeaderRoutes,
	employeeRoutes EmployeeRoutes,
	tarikdanaRoutes TarikDanaRoutes,
	invoiceheaderRoutes InvoiceHeaderRoutes,
	invoicedetailRoutes InvoiceDetailRoutes,
) Routes {
	return Routes{
		pprofRoutes,
		swaggerRoutes,
		publicRoutes,
		userRoutes,
		roleRoutes,
		menuRoutes,
		companyRoutes,
		branchRoutes,
		departmentRoutes,
		accountRoutes,
		costcentreRoutes,
		saldoRoutes,
		trxRoutes,
		kasbonRoutes,
		bkkdetailRoutes,
		bkkHeaderRoutes,
		employeeRoutes,
		tarikdanaRoutes,
		invoiceheaderRoutes,
		invoicedetailRoutes,
	}
}

// Setup all the route
func (a Routes) Setup() {
	for _, route := range a {
		route.Setup()
	}
}
