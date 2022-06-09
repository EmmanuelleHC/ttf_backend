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

type CompanyController struct {
	logger         lib.Logger
	companyService services.CompanyService
}

// NewCompanyController creates new company controller
func NewCompanyController(
	logger lib.Logger,
	companyService services.CompanyService,
) CompanyController {
	return CompanyController{
		logger:         logger,
		companyService: companyService,
	}
}

// @tags Company
// @summary Company Query
// @produce application/json
// @param data query models.CompanyQueryParam true "CompanyQueryParam"
// @success 200 {object} echox.Response{data=models.CompanyQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/companys [get]
func (a CompanyController) Query(ctx echo.Context) error {
	param := new(models.CompanyQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.companyService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags Company
// @summary Company Get All
// @produce application/json
// @param data query models.CompanyQueryParam true "CompanyQueryParam"
// @success 200 {object} echox.Response{data=models.Companys} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/companys [get]
func (a CompanyController) GetAll(ctx echo.Context) error {
	qr, err := a.companyService.Query(&models.CompanyQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags Company
// @summary Company Get By ID
// @produce application/json
// @param id path int true "company id"
// @success 200 {object} echox.Response{data=models.Company} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/companys/{id} [get]
func (a CompanyController) Get(ctx echo.Context) error {
	company, err := a.companyService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: company}.JSON(ctx)
}

// @tags Company
// @summary Company Create
// @produce application/json
// @param data body models.Company true "Company"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/companys [post]
func (a CompanyController) Create(ctx echo.Context) error {
	company := new(models.Company)
	if err := ctx.Bind(company); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	company.CreatedBy = claims.Username

	id, err := a.companyService.WithTrx(trxHandle).Create(company)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags Company
// @summary Company Update By ID
// @produce application/json
// @param id path int true "company id"
// @param data body models.Company true "Company"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/companys/{id} [put]
func (a CompanyController) Update(ctx echo.Context) error {
	company := new(models.Company)
	if err := ctx.Bind(company); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	company.UpdateBy = claims.Username

	if err := a.companyService.WithTrx(trxHandle).Update(ctx.Param("id"), company); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Company
// @summary Company Delete By ID
// @produce application/json
// @param id path int true "company id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/companys/{id} [delete]
func (a CompanyController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.companyService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Company
// @summary Company Enable By ID
// @produce application/json
// @param id path int true "company id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/companys/{id}/enable [patch]
func (a CompanyController) Enable(ctx echo.Context) error {
	if err := a.companyService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Company
// @summary Company Disable By ID
// @produce application/json
// @param id path int true "company id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/companys/{id}/disable [patch]
func (a CompanyController) Disable(ctx echo.Context) error {
	if err := a.companyService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}
