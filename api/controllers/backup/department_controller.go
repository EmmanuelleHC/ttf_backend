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

type DepartmentController struct {
	logger            lib.Logger
	departmentService services.DepartmentService
}

// NewDepartmentController creates new department controller
func NewDepartmentController(
	logger lib.Logger,
	departmentService services.DepartmentService,
) DepartmentController {
	return DepartmentController{
		logger:            logger,
		departmentService: departmentService,
	}
}

// @tags Department
// @summary Department Query
// @produce application/json
// @param data query models.DepartmentQueryParam true "DepartmentQueryParam"
// @success 200 {object} echox.Response{data=models.DepartmentQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/departments [get]
func (a DepartmentController) Query(ctx echo.Context) error {
	param := new(models.DepartmentQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.departmentService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags Department
// @summary Department Get All
// @produce application/json
// @param data query models.DepartmentQueryParam true "DepartmentQueryParam"
// @success 200 {object} echox.Response{data=models.Departments} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/departments [get]
func (a DepartmentController) GetAll(ctx echo.Context) error {
	qr, err := a.departmentService.Query(&models.DepartmentQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags Department
// @summary Department Get By ID
// @produce application/json
// @param id path int true "department id"
// @success 200 {object} echox.Response{data=models.Department} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/departments/{id} [get]
func (a DepartmentController) Get(ctx echo.Context) error {
	department, err := a.departmentService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: department}.JSON(ctx)
}

// @tags Department
// @summary Department Create
// @produce application/json
// @param data body models.Department true "Department"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/departments [post]
func (a DepartmentController) Create(ctx echo.Context) error {
	department := new(models.Department)
	if err := ctx.Bind(department); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	department.CreatedBy = claims.Username

	id, err := a.departmentService.WithTrx(trxHandle).Create(department)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags Department
// @summary Department Update By ID
// @produce application/json
// @param id path int true "department id"
// @param data body models.Department true "Department"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/departments/{id} [put]
func (a DepartmentController) Update(ctx echo.Context) error {
	department := new(models.Department)
	if err := ctx.Bind(department); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	department.UpdateBy = claims.Username

	if err := a.departmentService.WithTrx(trxHandle).Update(ctx.Param("id"), department); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Department
// @summary Department Delete By ID
// @produce application/json
// @param id path int true "department id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/departments/{id} [delete]
func (a DepartmentController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.departmentService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Department
// @summary Department Enable By ID
// @produce application/json
// @param id path int true "department id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/departments/{id}/enable [patch]
func (a DepartmentController) Enable(ctx echo.Context) error {
	if err := a.departmentService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Department
// @summary Department Disable By ID
// @produce application/json
// @param id path int true "department id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/departments/{id}/disable [patch]
func (a DepartmentController) Disable(ctx echo.Context) error {
	if err := a.departmentService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}
