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

type TarikDanaController struct {
	logger           lib.Logger
	tarikdanaService services.TarikDanaService
}

// NewTarikDanaController creates new tarikdana controller
func NewTarikDanaController(
	logger lib.Logger,
	tarikdanaService services.TarikDanaService,
) TarikDanaController {
	return TarikDanaController{
		logger:           logger,
		tarikdanaService: tarikdanaService,
	}
}

// @tags TarikDana
// @summary TarikDana Query
// @produce application/json
// @param data query models.TarikDanaQueryParam true "TarikDanaQueryParam"
// @success 200 {object} echox.Response{data=models.TarikDanaQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/tarikdanas [get]
func (a TarikDanaController) Query(ctx echo.Context) error {
	param := new(models.TarikDanaQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.tarikdanaService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags TarikDana
// @summary TarikDana Get All
// @produce application/json
// @param data query models.TarikDanaQueryParam true "TarikDanaQueryParam"
// @success 200 {object} echox.Response{data=models.TarikDanas} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/tarikdanas [get]
func (a TarikDanaController) GetAll(ctx echo.Context) error {
	qr, err := a.tarikdanaService.Query(&models.TarikDanaQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags TarikDana
// @summary TarikDana Get By ID
// @produce application/json
// @param id path int true "tarikdana id"
// @success 200 {object} echox.Response{data=models.TarikDana} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/tarikdanas/{id} [get]
func (a TarikDanaController) Get(ctx echo.Context) error {
	tarikdana, err := a.tarikdanaService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: tarikdana}.JSON(ctx)
}

// @tags TarikDana
// @summary TarikDana Create
// @produce application/json
// @param data body models.TarikDana true "TarikDana"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/tarikdanas [post]
func (a TarikDanaController) Create(ctx echo.Context) error {
	tarikdana := new(models.TarikDana)
	if err := ctx.Bind(tarikdana); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	tarikdana.CreatedBy = claims.Username

	id, err := a.tarikdanaService.WithTrx(trxHandle).Create(tarikdana)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags TarikDana
// @summary TarikDana Update By ID
// @produce application/json
// @param id path int true "tarikdana id"
// @param data body models.TarikDana true "TarikDana"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/tarikdanas/{id} [put]
func (a TarikDanaController) Update(ctx echo.Context) error {
	tarikdana := new(models.TarikDana)
	if err := ctx.Bind(tarikdana); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	tarikdana.UpdateBy = claims.Username

	if err := a.tarikdanaService.WithTrx(trxHandle).Update(ctx.Param("id"), tarikdana); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags TarikDana
// @summary TarikDana Delete By ID
// @produce application/json
// @param id path int true "tarikdana id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/tarikdanas/{id} [delete]
func (a TarikDanaController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.tarikdanaService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags TarikDana
// @summary TarikDana Enable By ID
// @produce application/json
// @param id path int true "tarikdana id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/tarikdanas/{id}/enable [patch]
func (a TarikDanaController) Enable(ctx echo.Context) error {
	if err := a.tarikdanaService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags TarikDana
// @summary TarikDana Disable By ID
// @produce application/json
// @param id path int true "tarikdana id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/tarikdanas/{id}/disable [patch]
func (a TarikDanaController) Disable(ctx echo.Context) error {
	if err := a.tarikdanaService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags TarikDana
// @summary TarikDana GetFile
// @produce application/json
// @param id path int true "tarikdana file id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/tarikdanas/upload [delete]
func (a TarikDanaController) GetFile(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	img := "./upload/dana/" + claims.Username + "/" + ctx.Param("id")

	return ctx.File(img)
}

// @tags TarikDana
// @summary TarikDana UploadFile By ID
// @produce application/json
// @param id path int true "tarikdana file"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/tarikdanas/upload [post]
func (a TarikDanaController) UploadFile(ctx echo.Context) error {
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

	if _, err := os.Stat("./upload/dana/" + claims.Username); os.IsNotExist(err) {
		err := os.Mkdir("./upload/dana/"+claims.Username, os.ModePerm)
		if err != nil {
			return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
		}
	}
	// Destination
	img := "./upload/dana/" + claims.Username + "/" + claims.Username + "-" + file.Filename
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

// @tags TarikDana
// @summary TarikDana RemoveFile
// @produce application/json
// @param id path int true "tarikdana file id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/tarikdanas/upload [delete]
func (a TarikDanaController) RemoveFile(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)

	// Delete
	img := "./upload/dana/" + claims.Username + "/" + claims.Username + "-" + ctx.Param("id")
	err := os.Remove(img)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Message: img}.JSON(ctx)
}
