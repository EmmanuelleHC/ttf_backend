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

type BKKDetailController struct {
	logger           lib.Logger
	bkkdetailService services.BKKDetailService
}

// NewBKKDetailController creates new bkkdetail controller
func NewBKKDetailController(
	logger lib.Logger,
	bkkdetailService services.BKKDetailService,
) BKKDetailController {
	return BKKDetailController{
		logger:           logger,
		bkkdetailService: bkkdetailService,
	}
}

// @tags BKKDetail
// @summary BKKDetail Query
// @produce application/json
// @param data query models.BKKDetailQueryParam true "BKKDetailQueryParam"
// @success 200 {object} echox.Response{data=models.BKKDetailQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkdetails [get]
func (a BKKDetailController) Query(ctx echo.Context) error {
	param := new(models.BKKDetailQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.bkkdetailService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags BKKDetail
// @summary BKKDetail Get All
// @produce application/json
// @param data query models.BKKDetailQueryParam true "BKKDetailQueryParam"
// @success 200 {object} echox.Response{data=models.BKKDetails} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkdetails [get]
func (a BKKDetailController) GetAll(ctx echo.Context) error {
	qr, err := a.bkkdetailService.Query(&models.BKKDetailQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags BKKDetail
// @summary BKKDetail Get By ID
// @produce application/json
// @param id path int true "bkkdetail id"
// @success 200 {object} echox.Response{data=models.BKKDetail} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkdetails/{id} [get]
func (a BKKDetailController) Get(ctx echo.Context) error {
	bkkdetail, err := a.bkkdetailService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: bkkdetail}.JSON(ctx)
}

// @tags BKKDetail
// @summary BKKDetail Create
// @produce application/json
// @param data body models.BKKDetail true "BKKDetail"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkdetails [post]
func (a BKKDetailController) Create(ctx echo.Context) error {
	bkkdetail := new(models.BKKDetail)
	if err := ctx.Bind(bkkdetail); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	bkkdetail.CreatedBy = claims.Username

	id, err := a.bkkdetailService.WithTrx(trxHandle).Create(bkkdetail)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags BKKDetail
// @summary BKKDetail Update By ID
// @produce application/json
// @param id path int true "bkkdetail id"
// @param data body models.BKKDetail true "BKKDetail"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkdetails/{id} [put]
func (a BKKDetailController) Update(ctx echo.Context) error {
	bkkdetail := new(models.BKKDetail)
	if err := ctx.Bind(bkkdetail); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	bkkdetail.UpdateBy = claims.Username

	if err := a.bkkdetailService.WithTrx(trxHandle).Update(ctx.Param("id"), bkkdetail); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags BKKDetail
// @summary BKKDetail Delete By ID
// @produce application/json
// @param id path int true "bkkdetail id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkdetails/{id} [delete]
func (a BKKDetailController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.bkkdetailService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags BKKDetail
// @summary BKKDetail Enable By ID
// @produce application/json
// @param id path int true "bkkdetail id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkdetails/{id}/enable [patch]
func (a BKKDetailController) Enable(ctx echo.Context) error {
	if err := a.bkkdetailService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags BKKDetail
// @summary BKKDetail Disable By ID
// @produce application/json
// @param id path int true "bkkdetail id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkdetails/{id}/disable [patch]
func (a BKKDetailController) Disable(ctx echo.Context) error {
	if err := a.bkkdetailService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags BKKDetail
// @summary BKKDetail GetFile
// @produce application/json
// @param id path int true "bkkdetail file id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkdetails/upload [delete]
func (a BKKDetailController) GetFile(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)

	img := "./upload/" + claims.Username + "/" + ctx.Param("id")
	return ctx.File(img)
}

// @tags BKKDetail
// @summary BKKDetail UploadFile By ID
// @produce application/json
// @param id path int true "bkkdetail file"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkdetails/upload [post]
func (a BKKDetailController) UploadFile(ctx echo.Context) error {
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

// @tags BKKDetail
// @summary BKKDetail RemoveFile
// @produce application/json
// @param id path int true "bkkdetail file id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/bkkdetails/upload [delete]
func (a BKKDetailController) RemoveFile(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)

	// Delete
	img := "./upload/" + claims.Username + "/" + claims.Username + "-" + ctx.Param("id")
	err := os.Remove(img)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Message: img}.JSON(ctx)
}
