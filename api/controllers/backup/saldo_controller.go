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

type SaldoController struct {
	logger       lib.Logger
	saldoService services.SaldoService
}

// NewSaldoController creates new saldo controller
func NewSaldoController(
	logger lib.Logger,
	saldoService services.SaldoService,
) SaldoController {
	return SaldoController{
		logger:       logger,
		saldoService: saldoService,
	}
}

// @tags Saldo
// @summary Saldo Query
// @produce application/json
// @param data query models.SaldoQueryParam true "SaldoQueryParam"
// @success 200 {object} echox.Response{data=models.SaldoQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/saldos [get]
func (a SaldoController) Query(ctx echo.Context) error {
	param := new(models.SaldoQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.saldoService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags Saldo
// @summary Saldo Get All
// @produce application/json
// @param data query models.SaldoQueryParam true "SaldoQueryParam"
// @success 200 {object} echox.Response{data=models.Saldos} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/saldos [get]
func (a SaldoController) GetAll(ctx echo.Context) error {
	qr, err := a.saldoService.Query(&models.SaldoQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags Saldo
// @summary Saldo Get By ID
// @produce application/json
// @param id path int true "saldo id"
// @success 200 {object} echox.Response{data=models.Saldo} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/saldos/{id} [get]
func (a SaldoController) Get(ctx echo.Context) error {
	saldo, err := a.saldoService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: saldo}.JSON(ctx)
}

// @tags Saldo
// @summary Saldo Get By User
// @produce application/json
// @success 200 {object} echox.Response{data=models.Saldo} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/saldos/user [get]
func (a SaldoController) GetByUser(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	saldo, err := a.saldoService.GetByUser(claims.Username)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: saldo}.JSON(ctx)
}

// @tags Saldo
// @summary Saldo Create
// @produce application/json
// @param data body models.Saldo true "Saldo"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/saldos [post]
func (a SaldoController) Create(ctx echo.Context) error {
	saldo := new(models.Saldo)
	if err := ctx.Bind(saldo); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	saldo.CreatedBy = claims.Username

	id, err := a.saldoService.WithTrx(trxHandle).CreateAndMonth(saldo)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags Saldo
// @summary Saldo Update By ID
// @produce application/json
// @param id path int true "saldo id"
// @param data body models.Saldo true "Saldo"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/saldos/{id} [put]
func (a SaldoController) Update(ctx echo.Context) error {
	saldo := new(models.Saldo)
	if err := ctx.Bind(saldo); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	saldo.UpdateBy = claims.Username

	if err := a.saldoService.WithTrx(trxHandle).Update(ctx.Param("id"), saldo); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Saldo
// @summary Saldo Delete By ID
// @produce application/json
// @param id path int true "saldo id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/saldos/{id} [delete]
func (a SaldoController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.saldoService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Saldo
// @summary Saldo Enable By ID
// @produce application/json
// @param id path int true "saldo id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/saldos/{id}/enable [patch]
func (a SaldoController) Enable(ctx echo.Context) error {
	if err := a.saldoService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Saldo
// @summary Saldo Disable By ID
// @produce application/json
// @param id path int true "saldo id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/saldos/{id}/disable [patch]
func (a SaldoController) Disable(ctx echo.Context) error {
	if err := a.saldoService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}
