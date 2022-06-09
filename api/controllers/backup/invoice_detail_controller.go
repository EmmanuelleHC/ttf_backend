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

type InvoiceDetailController struct {
	logger               lib.Logger
	invoicedetailService services.InvoiceDetailService
}

// NewInvoiceDetailController creates new invoicedetail controller
func NewInvoiceDetailController(
	logger lib.Logger,
	invoicedetailService services.InvoiceDetailService,
) InvoiceDetailController {
	return InvoiceDetailController{
		logger:               logger,
		invoicedetailService: invoicedetailService,
	}
}

// @tags InvoiceDetail
// @summary InvoiceDetail Query
// @produce application/json
// @param data query models.InvoiceDetailQueryParam true "InvoiceDetailQueryParam"
// @success 200 {object} echox.Response{data=models.InvoiceDetailQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoicedetails [get]
func (a InvoiceDetailController) Query(ctx echo.Context) error {
	param := new(models.InvoiceDetailQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.invoicedetailService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags InvoiceDetail
// @summary InvoiceDetail QueryBkk
// @produce application/json
// @param data query models.InvoiceDetailQueryParam true "InvoiceDetailQueryParam"
// @success 200 {object} echox.Response{data=models.InvoiceDetailQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoicedetails/bkk [get]
func (a InvoiceDetailController) QueryBkk(ctx echo.Context) error {
	param := new(models.InvoiceDetailQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.invoicedetailService.QueryBkk(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags InvoiceDetail
// @summary InvoiceDetail Get All
// @produce application/json
// @param data query models.InvoiceDetailQueryParam true "InvoiceDetailQueryParam"
// @success 200 {object} echox.Response{data=models.InvoiceDetails} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoicedetails [get]
func (a InvoiceDetailController) GetAll(ctx echo.Context) error {
	qr, err := a.invoicedetailService.Query(&models.InvoiceDetailQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags InvoiceDetail
// @summary InvoiceDetail Get By ID
// @produce application/json
// @param id path int true "invoicedetail id"
// @success 200 {object} echox.Response{data=models.InvoiceDetail} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoicedetails/{id} [get]
func (a InvoiceDetailController) Get(ctx echo.Context) error {
	invoicedetail, err := a.invoicedetailService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: invoicedetail}.JSON(ctx)
}

// @tags InvoiceDetail
// @summary InvoiceDetail Create
// @produce application/json
// @param data body models.InvoiceDetail true "InvoiceDetail"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoicedetails [post]
func (a InvoiceDetailController) Create(ctx echo.Context) error {
	invoicedetail := new(models.InvoiceDetail)
	if err := ctx.Bind(invoicedetail); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	invoicedetail.CreatedBy = claims.Username

	id, err := a.invoicedetailService.WithTrx(trxHandle).Create(invoicedetail)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags InvoiceDetail
// @summary InvoiceDetail Update By ID
// @produce application/json
// @param id path int true "invoicedetail id"
// @param data body models.InvoiceDetail true "InvoiceDetail"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoicedetails/{id} [put]
func (a InvoiceDetailController) Update(ctx echo.Context) error {
	invoicedetail := new(models.InvoiceDetail)
	if err := ctx.Bind(invoicedetail); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	invoicedetail.UpdateBy = claims.Username

	if err := a.invoicedetailService.WithTrx(trxHandle).Update(ctx.Param("id"), invoicedetail); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags InvoiceDetail
// @summary InvoiceDetail Delete By ID
// @produce application/json
// @param id path int true "invoicedetail id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoicedetails/{id} [delete]
func (a InvoiceDetailController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.invoicedetailService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags InvoiceDetail
// @summary InvoiceDetail Enable By ID
// @produce application/json
// @param id path int true "invoicedetail id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoicedetails/{id}/enable [patch]
func (a InvoiceDetailController) Enable(ctx echo.Context) error {
	if err := a.invoicedetailService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags InvoiceDetail
// @summary InvoiceDetail Disable By ID
// @produce application/json
// @param id path int true "invoicedetail id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoicedetails/{id}/disable [patch]
func (a InvoiceDetailController) Disable(ctx echo.Context) error {
	if err := a.invoicedetailService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags InvoiceDetail
// @summary InvoiceDetail GetFile
// @produce application/json
// @param id path int true "invoicedetail file id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoicedetails/upload [delete]
func (a InvoiceDetailController) GetFile(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)

	img := "./upload/" + claims.Username + "/" + ctx.Param("id")
	return ctx.File(img)
}

// @tags InvoiceDetail
// @summary InvoiceDetail UploadFile By ID
// @produce application/json
// @param id path int true "invoicedetail file"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoicedetails/upload [post]
func (a InvoiceDetailController) UploadFile(ctx echo.Context) error {
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

	if _, err := os.Stat("./upload/" + claims.Username); os.IsNotExist(err) {
		err := os.Mkdir("./upload/"+claims.Username, os.ModePerm)
		if err != nil {
			return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
		}
	}
	// Destination
	img := "./upload/" + claims.Username + "/" + claims.Username + "-" + file.Filename
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

// @tags InvoiceDetail
// @summary InvoiceDetail RemoveFile
// @produce application/json
// @param id path int true "invoicedetail file id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/invoicedetails/upload [delete]
func (a InvoiceDetailController) RemoveFile(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)

	// Delete
	img := "./upload/" + claims.Username + "/" + claims.Username + "-" + ctx.Param("id")
	err := os.Remove(img)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Message: img}.JSON(ctx)
}
