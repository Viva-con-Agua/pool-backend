package token

import (
	"bytes"
	"net/http"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/labstack/echo/v4"
)

type ReceiptFileHandler struct {
	vcago.Handler
}

var ReceiptFile = &ReceiptFileHandler{*vcago.NewHandler("receipt")}

func (i *ReceiptFileHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Upload, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.GET("/zip/:id", i.GetZipByID, accessCookie)
	group.DELETE("/:id", i.DeleteByID, accessCookie)

}

func (i *ReceiptFileHandler) Upload(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	body := new(models.ReceiptFileCreate)
	if err = c.BindFormDataAndValidate("body", body); err != nil {
		return
	}
	file := new(vmod.File)
	if file, err = c.BindFormDataFile("image"); err != nil {
		return
	}
	var result *models.ReceiptFile
	if result, err = dao.ReceiptFileCreate(c.Ctx(), body, file, token); err != nil {
		return
	}
	return c.Created(result)
}

func (i *ReceiptFileHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	body := new(vmod.IDParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	var result []byte
	if result, err = dao.ReceiptFileGetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Stream(http.StatusOK, "image/png", bytes.NewReader(result))
}

func (i *ReceiptFileHandler) GetZipByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	body := new(vmod.IDParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	var result []byte
	if result, err = dao.ReceiptFileZipGetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Stream(http.StatusOK, "application/zip", bytes.NewReader(result))
}

func (i *ReceiptFileHandler) DeleteByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	body := new(vmod.IDParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	var result *vmod.DeletedResponse
	if result, err = dao.ReceiptFileDeleteByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Deleted(result)
}
