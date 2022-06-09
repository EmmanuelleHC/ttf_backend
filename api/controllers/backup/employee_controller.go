package controllers

import (
	"net/http"

	"github.com/Aguztinus/petty-cash-backend/api/services"
	"github.com/Aguztinus/petty-cash-backend/constants"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
	"github.com/Aguztinus/petty-cash-backend/pkg/echox"
	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type EmployeeController struct {
	logger          lib.Logger
	employeeService services.EmployeeService
}

// NewEmployeeController creates new employee controller
func NewEmployeeController(
	logger lib.Logger,
	employeeService services.EmployeeService,
) EmployeeController {
	return EmployeeController{
		logger:          logger,
		employeeService: employeeService,
	}
}

// @tags Employee
// @summary Employee Query
// @produce application/json
// @param data query models.EmployeeQueryParam true "EmployeeQueryParam"
// @success 200 {object} echox.Response{data=models.EmployeeQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/employees [get]
func (a EmployeeController) Query(ctx echo.Context) error {
	param := new(models.EmployeeQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.employeeService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags Employee
// @summary Employee Get All
// @produce application/json
// @param data query models.EmployeeQueryParam true "EmployeeQueryParam"
// @success 200 {object} echox.Response{data=models.Employees} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/employees [get]
func (a EmployeeController) GetAll(ctx echo.Context) error {
	qr, err := a.employeeService.Query(&models.EmployeeQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags Employee
// @summary Employee Get By ID
// @produce application/json
// @param id path int true "employee id"
// @success 200 {object} echox.Response{data=models.Employee} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/employees/{id} [get]
func (a EmployeeController) Get(ctx echo.Context) error {
	employee, err := a.employeeService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: employee}.JSON(ctx)
}

// @tags Employee
// @summary Employee Create
// @produce application/json
// @param data body models.Employee true "Employee"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/employees [post]
func (a EmployeeController) Create(ctx echo.Context) error {
	employee := new(models.Employee)
	if err := ctx.Bind(employee); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	employee.CreatedBy = claims.Username

	id, err := a.employeeService.WithTrx(trxHandle).Create(employee)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags Employee
// @summary Employee Update By ID
// @produce application/json
// @param id path int true "employee id"
// @param data body models.Employee true "Employee"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/employees/{id} [put]
func (a EmployeeController) Update(ctx echo.Context) error {
	employee := new(models.Employee)
	if err := ctx.Bind(employee); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	employee.UpdateBy = claims.Username

	if err := a.employeeService.WithTrx(trxHandle).Update(ctx.Param("id"), employee); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Employee
// @summary Employee Delete By ID
// @produce application/json
// @param id path int true "employee id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/employees/{id} [delete]
func (a EmployeeController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.employeeService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Employee
// @summary Employee Enable By ID
// @produce application/json
// @param id path int true "employee id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/employees/{id}/enable [patch]
func (a EmployeeController) Enable(ctx echo.Context) error {
	if err := a.employeeService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Employee
// @summary Employee Disable By ID
// @produce application/json
// @param id path int true "employee id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/employees/{id}/disable [patch]
func (a EmployeeController) Disable(ctx echo.Context) error {
	if err := a.employeeService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}
