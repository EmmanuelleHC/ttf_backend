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

type AccountController struct {
	logger         lib.Logger
	accountService services.AccountService
}

// NewAccountController creates new account controller
func NewAccountController(
	logger lib.Logger,
	accountService services.AccountService,
) AccountController {
	return AccountController{
		logger:         logger,
		accountService: accountService,
	}
}

// @tags Account
// @summary Account Query
// @produce application/json
// @param data query models.AccountQueryParam true "AccountQueryParam"
// @success 200 {object} echox.Response{data=models.AccountQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/accounts [get]
func (a AccountController) Query(ctx echo.Context) error {
	param := new(models.AccountQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.accountService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags Account
// @summary Account Get All
// @produce application/json
// @param data query models.AccountQueryParam true "AccountQueryParam"
// @success 200 {object} echox.Response{data=models.Accounts} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/accounts [get]
func (a AccountController) GetAll(ctx echo.Context) error {
	qr, err := a.accountService.Query(&models.AccountQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags Account
// @summary Account Get By ID
// @produce application/json
// @param id path int true "account id"
// @success 200 {object} echox.Response{data=models.Account} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/accounts/{id} [get]
func (a AccountController) Get(ctx echo.Context) error {
	account, err := a.accountService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: account}.JSON(ctx)
}

// @tags Account
// @summary Account Create
// @produce application/json
// @param data body models.Account true "Account"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/accounts [post]
func (a AccountController) Create(ctx echo.Context) error {
	account := new(models.Account)
	if err := ctx.Bind(account); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	account.CreatedBy = claims.Username

	id, err := a.accountService.WithTrx(trxHandle).Create(account)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags Account
// @summary Account Update By ID
// @produce application/json
// @param id path int true "account id"
// @param data body models.Account true "Account"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/accounts/{id} [put]
func (a AccountController) Update(ctx echo.Context) error {
	account := new(models.Account)
	if err := ctx.Bind(account); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	account.UpdateBy = claims.Username

	if err := a.accountService.WithTrx(trxHandle).Update(ctx.Param("id"), account); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Account
// @summary Account Delete By ID
// @produce application/json
// @param id path int true "account id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/accounts/{id} [delete]
func (a AccountController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.accountService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Account
// @summary Account Enable By ID
// @produce application/json
// @param id path int true "account id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/accounts/{id}/enable [patch]
func (a AccountController) Enable(ctx echo.Context) error {
	if err := a.accountService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Account
// @summary Account Disable By ID
// @produce application/json
// @param id path int true "account id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/accounts/{id}/disable [patch]
func (a AccountController) Disable(ctx echo.Context) error {
	if err := a.accountService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}
