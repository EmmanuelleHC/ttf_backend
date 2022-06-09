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

type KasbonController struct {
	logger        lib.Logger
	kasbonService services.KasbonService
}

// NewKasbonController creates new kasbon controller
func NewKasbonController(
	logger lib.Logger,
	kasbonService services.KasbonService,
) KasbonController {
	return KasbonController{
		logger:        logger,
		kasbonService: kasbonService,
	}
}

// @tags Kasbon
// @summary Kasbon Query
// @produce application/json
// @param data query models.KasbonQueryParam true "KasbonQueryParam"
// @success 200 {object} echox.Response{data=models.KasbonQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/kasbons [get]
func (a KasbonController) Query(ctx echo.Context) error {
	param := new(models.KasbonQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.kasbonService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags Kasbon
// @summary Kasbon Get All
// @produce application/json
// @param data query models.KasbonQueryParam true "KasbonQueryParam"
// @success 200 {object} echox.Response{data=models.Kasbons} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/kasbons [get]
func (a KasbonController) GetAll(ctx echo.Context) error {
	qr, err := a.kasbonService.Query(&models.KasbonQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags Kasbon
// @summary Kasbon Get By ID
// @produce application/json
// @param id path int true "kasbon id"
// @success 200 {object} echox.Response{data=models.Kasbon} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/kasbons/{id} [get]
func (a KasbonController) Get(ctx echo.Context) error {
	kasbon, err := a.kasbonService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: kasbon}.JSON(ctx)
}

// @tags Kasbon
// @summary Kasbon Create
// @produce application/json
// @param data body models.Kasbon true "Kasbon"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/kasbons [post]
func (a KasbonController) Create(ctx echo.Context) error {
	kasbon := new(models.Kasbon)
	if err := ctx.Bind(kasbon); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	kasbon.CreatedBy = claims.Username

	id, err := a.kasbonService.WithTrx(trxHandle).Create(kasbon)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags Kasbon
// @summary Kasbon Update By ID
// @produce application/json
// @param id path int true "kasbon id"
// @param data body models.Kasbon true "Kasbon"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/kasbons/{id} [put]
func (a KasbonController) Update(ctx echo.Context) error {
	kasbon := new(models.Kasbon)
	if err := ctx.Bind(kasbon); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	kasbon.UpdateBy = claims.Username

	if err := a.kasbonService.WithTrx(trxHandle).Update(ctx.Param("id"), kasbon); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Kasbon
// @summary Kasbon Delete By ID
// @produce application/json
// @param id path int true "kasbon id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/kasbons/{id} [delete]
func (a KasbonController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.kasbonService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Kasbon
// @summary Kasbon GetFile
// @produce application/json
// @param id path int true "kasbon file id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/kasbons/upload [delete]
func (a KasbonController) GetFile(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	img := "./upload/" + claims.Username + "/" + ctx.Param("id")

	return ctx.File(img)
}

// @tags Kasbon
// @summary Kasbon UploadFile By ID
// @produce application/json
// @param id path int true "kasbon file"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/kasbons/upload [post]
func (a KasbonController) UploadFile(ctx echo.Context) error {
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

// @tags Kasbon
// @summary Kasbon RemoveFile
// @produce application/json
// @param id path int true "kasbon file id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/kasbons/upload [delete]
func (a KasbonController) RemoveFile(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)

	// Delete
	img := "./upload/" + claims.Username + "/" + claims.Username + "-" + ctx.Param("id")
	err := os.Remove(img)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Message: img}.JSON(ctx)
}
