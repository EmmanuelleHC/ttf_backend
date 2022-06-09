package controllers

import (
	"io"
	"net/http"
	"os"

	"github.com/Aguztinus/petty-cash-backend/api/services"
	"github.com/Aguztinus/petty-cash-backend/constants"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
	"github.com/Aguztinus/petty-cash-backend/pkg/echox"
	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type InvoiceHeaderController struct {
	logger               lib.Logger
	invoiceheaderService services.InvoiceHeaderService
}

// NewInvoiceHeaderController creates new invoiceheader controller
func NewInvoiceHeaderController(
	logger lib.Logger,
	invoiceheaderService services.InvoiceHeaderService,
) InvoiceHeaderController {
	return InvoiceHeaderController{
		logger:               logger,
		invoiceheaderService: invoiceheaderService,
	}
}

// @tags InvoiceHeader
// @summary InvoiceHeader Query
// @produce application/json
// @param data query models.InvoiceHeaderQueryParam true "InvoiceHeaderQueryParam"
// @success 200 {object} echox.Response{data=models.InvoiceHeaderQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoiceheaders [get]
func (a InvoiceHeaderController) Query(ctx echo.Context) error {
	param := new(models.InvoiceHeaderQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.invoiceheaderService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags InvoiceHeader
// @summary InvoiceHeader Get All
// @produce application/json
// @param data query models.InvoiceHeaderQueryParam true "InvoiceHeaderQueryParam"
// @success 200 {object} echox.Response{data=models.InvoiceHeaders} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoiceheaders [get]
func (a InvoiceHeaderController) GetAll(ctx echo.Context) error {
	qr, err := a.invoiceheaderService.Query(&models.InvoiceHeaderQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags InvoiceHeader
// @summary InvoiceHeader Get By ID
// @produce application/json
// @param id path int true "invoiceheader id"
// @success 200 {object} echox.Response{data=models.InvoiceHeader} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoiceheaders/{id} [get]
func (a InvoiceHeaderController) Get(ctx echo.Context) error {
	invoiceheader, err := a.invoiceheaderService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: invoiceheader}.JSON(ctx)
}

// @tags InvoiceHeader
// @summary InvoiceHeader Create
// @produce application/json
// @param data body models.InvoiceHeader true "InvoiceHeader"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoiceheaders [post]
func (a InvoiceHeaderController) Create(ctx echo.Context) error {
	invoiceheader := new(models.InvoiceHeader)
	if err := ctx.Bind(invoiceheader); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	invoiceheader.CreatedBy = claims.Username

	id, err := a.invoiceheaderService.WithTrx(trxHandle).Create(invoiceheader)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags InvoiceHeader
// @summary InvoiceHeader Update By ID
// @produce application/json
// @param id path int true "invoiceheader id"
// @param data body models.InvoiceHeader true "InvoiceHeader"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoiceheaders/{id} [put]
func (a InvoiceHeaderController) Update(ctx echo.Context) error {
	invoiceheader := new(models.InvoiceHeader)
	if err := ctx.Bind(invoiceheader); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	invoiceheader.UpdateBy = claims.Username

	if err := a.invoiceheaderService.WithTrx(trxHandle).Update(ctx.Param("id"), invoiceheader); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags InvoiceHeader
// @summary InvoiceHeader Delete By ID
// @produce application/json
// @param id path int true "invoiceheader id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoiceheaders/{id} [delete]
func (a InvoiceHeaderController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.invoiceheaderService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags InvoiceHeader
// @summary InvoiceHeader Enable By ID
// @produce application/json
// @param id path int true "invoiceheader id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoiceheaders/{id}/enable [patch]
func (a InvoiceHeaderController) Enable(ctx echo.Context) error {
	if err := a.invoiceheaderService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags InvoiceHeader
// @summary InvoiceHeader Disable By ID
// @produce application/json
// @param id path int true "invoiceheader id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoiceheaders/{id}/disable [patch]
func (a InvoiceHeaderController) Disable(ctx echo.Context) error {
	if err := a.invoiceheaderService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags invoiceheader
// @summary invoiceheader StatusApprove By ID
// @produce application/json
// @param id path int true "invoiceheader id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoiceheader/{id}/statusapprove [patch]
func (a InvoiceHeaderController) StatusApprove(ctx echo.Context) error {
	ids := new(models.InvoiceHeaderApproveQueryParam)
	if err := ctx.Bind(ids); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	if err := a.invoiceheaderService.UpdateApprove(ids.IDs, 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags InvoiceHeader
// @summary InvoiceHeader StatusReject By ID
// @produce application/json
// @param id path int true "invoiceHeader id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoiceheader/{id}/statusreject [patch]
func (a InvoiceHeaderController) StatusReject(ctx echo.Context) error {
	ids := new(models.InvoiceHeaderApproveQueryParam)
	if err := ctx.Bind(ids); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	if err := a.invoiceheaderService.UpdateApprove(ids.IDs, 2); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags invoiceheader
// @summary invoiceheader StatusApproveFinal By ID
// @produce application/json
// @param id path int true "invoiceheader id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoiceheader/{id}/statusapprovefinal [patch]
func (a InvoiceHeaderController) StatusApproveFinal(ctx echo.Context) error {
	ids := new(models.InvoiceHeaderApproveQueryParam)
	if err := ctx.Bind(ids); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	if err := a.invoiceheaderService.UpdateApprove(ids.IDs, 3); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags InvoiceHeader
// @summary InvoiceHeader StatusRejectFinal By ID
// @produce application/json
// @param id path int true "invoiceHeader id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoiceheader/{id}/statusrejectfinal [patch]
func (a InvoiceHeaderController) StatusRejectFinal(ctx echo.Context) error {
	ids := new(models.InvoiceHeaderApproveQueryParam)
	if err := ctx.Bind(ids); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	if err := a.invoiceheaderService.UpdateApprove(ids.IDs, 4); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Invoice
// @summary Invoice GetFile
// @produce application/json
// @param id path int true "invoice file id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoices/upload [delete]
func (a InvoiceHeaderController) GetFile(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	img := "./upload/invoice/" + claims.Username + "/" + ctx.Param("id")

	return ctx.File(img)
}

// @tags Invoice
// @summary Invoice UploadFile By ID
// @produce application/json
// @param id path int true "invoice file"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoices/upload [post]
func (a InvoiceHeaderController) UploadFile(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)

	// Source
	file, err := ctx.FormFile("file")
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}
	src, err := file.Open()
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}
	defer src.Close()

	if _, err := os.Stat("./upload/invoice/" + claims.Username); os.IsNotExist(err) {
		err := os.Mkdir("./upload/invoice/"+claims.Username, os.ModePerm)
		if err != nil {
			return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
		}
	}
	// Destination
	img := "./upload/invoice/" + claims.Username + "/" + claims.Username + "-" + file.Filename
	dst, err := os.Create(img)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}
	return echox.Response{Code: http.StatusOK, Message: claims.Username + "-" + file.Filename}.JSON(ctx)
}

// @tags Invoice
// @summary Invoice RemoveFile
// @produce application/json
// @param id path int true "invoice file id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoices/upload [delete]
func (a InvoiceHeaderController) RemoveFile(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)

	// Delete
	img := "./upload/invoice/" + claims.Username + "/" + claims.Username + "-" + ctx.Param("id")
	err := os.Remove(img)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Message: img}.JSON(ctx)
}
