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

type BKKHeaderController struct {
	logger           lib.Logger
	bkkheaderService services.BKKHeaderService
	reportService    services.ReportService
}

// NewBKKHeaderController creates new bkkheader controller
func NewBKKHeaderController(
	logger lib.Logger,
	bkkheaderService services.BKKHeaderService,
	reportService services.ReportService,
) BKKHeaderController {
	return BKKHeaderController{
		logger:           logger,
		bkkheaderService: bkkheaderService,
		reportService:    reportService,
	}
}

// @tags BKKHeader
// @summary BKKHeader Query
// @produce application/json
// @param data query models.BKKHeaderQueryParam true "BKKHeaderQueryParam"
// @success 200 {object} echox.Response{data=models.BKKHeaderQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkheaders [get]
func (a BKKHeaderController) Query(ctx echo.Context) error {
	param := new(models.BKKHeaderQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.bkkheaderService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags BKKHeader
// @summary BKKHeader Query
// @produce application/json
// @param data query models.BKKHeaderQueryParam true "BKKHeaderQueryParam"
// @success 200 {object} echox.Response{data=models.BKKHeaderQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkheaders/approve [get]
func (a BKKHeaderController) GetApprove(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	param := new(models.BKKHeaderQueryParam)
	param.UserId = claims.ID
	qr, err := a.bkkheaderService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags BKKHeader
// @summary BKKHeader Get All
// @produce application/json
// @param data query models.BKKHeaderQueryParam true "BKKHeaderQueryParam"
// @success 200 {object} echox.Response{data=models.BKKHeaders} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkheaders [get]
func (a BKKHeaderController) GetAll(ctx echo.Context) error {
	qr, err := a.bkkheaderService.Query(&models.BKKHeaderQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags BKKHeader
// @summary BKKHeader Get By ID
// @produce application/json
// @param id path int true "bkkheader id"
// @success 200 {object} echox.Response{data=models.BKKHeader} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkheaders/{id} [get]
func (a BKKHeaderController) Get(ctx echo.Context) error {
	bkkheader, err := a.bkkheaderService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: bkkheader}.JSON(ctx)
}

// @tags BKKHeader
// @summary BKKHeader Create
// @produce application/json
// @param data body models.BKKHeader true "BKKHeader"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkheaders [post]
func (a BKKHeaderController) Create(ctx echo.Context) error {
	bkkheader := new(models.BKKHeader)
	if err := ctx.Bind(bkkheader); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	bkkheader.CreatedBy = claims.Username

	id, err := a.bkkheaderService.WithTrx(trxHandle).Create(bkkheader)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags BKKHeader
// @summary BKKHeader Update By ID
// @produce application/json
// @param id path int true "bkkheader id"
// @param data body models.BKKHeader true "BKKHeader"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkheaders/{id} [put]
func (a BKKHeaderController) Update(ctx echo.Context) error {
	bkkheader := new(models.BKKHeader)
	if err := ctx.Bind(bkkheader); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	bkkheader.UpdateBy = claims.Username

	if err := a.bkkheaderService.WithTrx(trxHandle).Update(ctx.Param("id"), bkkheader); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags BKKHeader
// @summary BKKHeader Delete By ID
// @produce application/json
// @param id path int true "bkkheader id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkheaders/{id} [delete]
func (a BKKHeaderController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.bkkheaderService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags BKKHeader
// @summary BKKHeader StatusPaid By ID
// @produce application/json
// @param id path int true "bkkheader id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkheaders/{id}/statuspaid [patch]
func (a BKKHeaderController) StatusPaid(ctx echo.Context) error {
	if err := a.bkkheaderService.UpdateStatus(ctx.Param("id"), "Paid", 0); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags BKKHeader
// @summary BKKHeader StatusApprove By ID
// @produce application/json
// @param id path int true "bkkheader id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkheaders/{id}/statusapprove [patch]
func (a BKKHeaderController) StatusApprove(ctx echo.Context) error {
	ids := new(models.BKKHeaderApproveQueryParam)
	if err := ctx.Bind(ids); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	if err := a.bkkheaderService.UpdateApprove(ids.IDs, 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags BKKHeader
// @summary BKKHeader StatusReject By ID
// @produce application/json
// @param id path int true "bkkheader id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkheaders/{id}/statusreject [patch]
func (a BKKHeaderController) StatusReject(ctx echo.Context) error {
	ids := new(models.BKKHeaderApproveQueryParam)
	if err := ctx.Bind(ids); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	if err := a.bkkheaderService.UpdateApprove(ids.IDs, 2); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}
