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

type TrxController struct {
	logger     lib.Logger
	trxService services.TrxService
}

// NewTrxController creates new trx controller
func NewTrxController(
	logger lib.Logger,
	trxService services.TrxService,
) TrxController {
	return TrxController{
		logger:     logger,
		trxService: trxService,
	}
}

// @tags Trx
// @summary Trx Query
// @produce application/json
// @param data query models.TrxQueryParam true "TrxQueryParam"
// @success 200 {object} echox.Response{data=models.TrxQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/trxs [get]
func (a TrxController) Query(ctx echo.Context) error {
	param := new(models.TrxQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.trxService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags Trx
// @summary Trx Get All
// @produce application/json
// @param data query models.TrxQueryParam true "TrxQueryParam"
// @success 200 {object} echox.Response{data=models.Trxs} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/trxs [get]
func (a TrxController) GetAll(ctx echo.Context) error {
	qr, err := a.trxService.Query(&models.TrxQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags Trx
// @summary Trx Get By ID
// @produce application/json
// @param id path int true "trx id"
// @success 200 {object} echox.Response{data=models.Trx} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/trxs/{id} [get]
func (a TrxController) Get(ctx echo.Context) error {
	trx, err := a.trxService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: trx}.JSON(ctx)
}

// @tags Trx
// @summary Trx Create
// @produce application/json
// @param data body models.Trx true "Trx"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/trxs [post]
func (a TrxController) Create(ctx echo.Context) error {
	trx := new(models.Trx)
	if err := ctx.Bind(trx); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	trx.CreatedBy = claims.Username

	id, err := a.trxService.WithTrx(trxHandle).Create(trx)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags Trx
// @summary Trx Update By ID
// @produce application/json
// @param id path int true "trx id"
// @param data body models.Trx true "Trx"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/trxs/{id} [put]
func (a TrxController) Update(ctx echo.Context) error {
	trx := new(models.Trx)
	if err := ctx.Bind(trx); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	trx.UpdateBy = claims.Username

	if err := a.trxService.WithTrx(trxHandle).Update(ctx.Param("id"), trx); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Trx
// @summary Trx Delete By ID
// @produce application/json
// @param id path int true "trx id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/trxs/{id} [delete]
func (a TrxController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.trxService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Trx
// @summary Trx Enable By ID
// @produce application/json
// @param id path int true "trx id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/trxs/{id}/enable [patch]
func (a TrxController) Enable(ctx echo.Context) error {
	if err := a.trxService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Trx
// @summary Trx Disable By ID
// @produce application/json
// @param id path int true "trx id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/trxs/{id}/disable [patch]
func (a TrxController) Disable(ctx echo.Context) error {
	if err := a.trxService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}
