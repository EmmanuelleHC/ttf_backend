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

type BranchController struct {
	logger        lib.Logger
	branchService services.BranchService
}

// NewBranchController creates new branch controller
func NewBranchController(
	logger lib.Logger,
	branchService services.BranchService,
) BranchController {
	return BranchController{
		logger:        logger,
		branchService: branchService,
	}
}

// @tags Branch
// @summary Branch Query
// @produce application/json
// @param data query models.BranchQueryParam true "BranchQueryParam"
// @success 200 {object} echox.Response{data=models.BranchQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/branchs [get]
func (a BranchController) Query(ctx echo.Context) error {
	param := new(models.BranchQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.branchService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags Branch
// @summary Branch Get All
// @produce application/json
// @param data query models.BranchQueryParam true "BranchQueryParam"
// @success 200 {object} echox.Response{data=models.Branchs} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/branchs [get]
func (a BranchController) GetAll(ctx echo.Context) error {
	qr, err := a.branchService.Query(&models.BranchQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags Branch
// @summary Branch Get By ID
// @produce application/json
// @param id path int true "branch id"
// @success 200 {object} echox.Response{data=models.Branch} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/branchs/{id} [get]
func (a BranchController) Get(ctx echo.Context) error {
	branch, err := a.branchService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: branch}.JSON(ctx)
}

// @tags Branch
// @summary Branch Create
// @produce application/json
// @param data body models.Branch true "Branch"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/branchs [post]
func (a BranchController) Create(ctx echo.Context) error {
	branch := new(models.Branch)
	if err := ctx.Bind(branch); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	branch.CreatedBy = claims.Username

	id, err := a.branchService.WithTrx(trxHandle).Create(branch)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags Branch
// @summary Branch Update By ID
// @produce application/json
// @param id path int true "branch id"
// @param data body models.Branch true "Branch"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/branchs/{id} [put]
func (a BranchController) Update(ctx echo.Context) error {
	branch := new(models.Branch)
	if err := ctx.Bind(branch); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	branch.UpdateBy = claims.Username

	if err := a.branchService.WithTrx(trxHandle).Update(ctx.Param("id"), branch); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Branch
// @summary Branch Delete By ID
// @produce application/json
// @param id path int true "branch id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/branchs/{id} [delete]
func (a BranchController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.branchService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Branch
// @summary Branch Enable By ID
// @produce application/json
// @param id path int true "branch id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/branchs/{id}/enable [patch]
func (a BranchController) Enable(ctx echo.Context) error {
	if err := a.branchService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Branch
// @summary Branch Disable By ID
// @produce application/json
// @param id path int true "branch id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/branchs/{id}/disable [patch]
func (a BranchController) Disable(ctx echo.Context) error {
	if err := a.branchService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}
