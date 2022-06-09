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

type CostCentreController struct {
	logger            lib.Logger
	costcentreService services.CostCentreService
}

// NewCostCentreController creates new costcentre controller
func NewCostCentreController(
	logger lib.Logger,
	costcentreService services.CostCentreService,
) CostCentreController {
	return CostCentreController{
		logger:            logger,
		costcentreService: costcentreService,
	}
}

// @tags CostCentre
// @summary CostCentre Query
// @produce application/json
// @param data query models.CostCentreQueryParam true "CostCentreQueryParam"
// @success 200 {object} echox.Response{data=models.CostCentreQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/costcentres [get]
func (a CostCentreController) Query(ctx echo.Context) error {
	param := new(models.CostCentreQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.costcentreService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags CostCentre
// @summary CostCentre Get All
// @produce application/json
// @param data query models.CostCentreQueryParam true "CostCentreQueryParam"
// @success 200 {object} echox.Response{data=models.CostCentres} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/costcentres [get]
func (a CostCentreController) GetAll(ctx echo.Context) error {
	qr, err := a.costcentreService.Query(&models.CostCentreQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags CostCentre
// @summary CostCentre Get By ID
// @produce application/json
// @param id path int true "costcentre id"
// @success 200 {object} echox.Response{data=models.CostCentre} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/costcentres/{id} [get]
func (a CostCentreController) Get(ctx echo.Context) error {
	costcentre, err := a.costcentreService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: costcentre}.JSON(ctx)
}

// @tags CostCentre
// @summary CostCentre Create
// @produce application/json
// @param data body models.CostCentre true "CostCentre"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/costcentres [post]
func (a CostCentreController) Create(ctx echo.Context) error {
	costcentre := new(models.CostCentre)
	if err := ctx.Bind(costcentre); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	costcentre.CreatedBy = claims.Username

	id, err := a.costcentreService.WithTrx(trxHandle).Create(costcentre)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags CostCentre
// @summary CostCentre Update By ID
// @produce application/json
// @param id path int true "costcentre id"
// @param data body models.CostCentre true "CostCentre"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/costcentres/{id} [put]
func (a CostCentreController) Update(ctx echo.Context) error {
	costcentre := new(models.CostCentre)
	if err := ctx.Bind(costcentre); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	costcentre.UpdateBy = claims.Username

	if err := a.costcentreService.WithTrx(trxHandle).Update(ctx.Param("id"), costcentre); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags CostCentre
// @summary CostCentre Delete By ID
// @produce application/json
// @param id path int true "costcentre id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/costcentres/{id} [delete]
func (a CostCentreController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.costcentreService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags CostCentre
// @summary CostCentre Enable By ID
// @produce application/json
// @param id path int true "costcentre id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/costcentres/{id}/enable [patch]
func (a CostCentreController) Enable(ctx echo.Context) error {
	if err := a.costcentreService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags CostCentre
// @summary CostCentre Disable By ID
// @produce application/json
// @param id path int true "costcentre id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/costcentres/{id}/disable [patch]
func (a CostCentreController) Disable(ctx echo.Context) error {
	if err := a.costcentreService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}
