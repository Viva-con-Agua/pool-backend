package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type ArtistHandler struct {
	vcago.Handler
}

var Artist = &ArtistHandler{*vcago.NewHandler("artist")}

func (i *ArtistHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)
}

func (i *ArtistHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ArtistCreate)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Artist)
	if result, err = dao.ArtistInsert(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Created(result)
}

func (i *ArtistHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ArtistQuery)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	result := new([]models.Artist)
	if result, err = dao.ArtistGet(c.Ctx(), body); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *ArtistHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ArtistParam)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	result := new(models.Artist)
	if result, err = dao.ArtistGetByID(c.Ctx(), body); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *ArtistHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ArtistUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Artist)
	if result, err = dao.ArtistUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Updated(result)
}

func (i *ArtistHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ArtistParam)
	if c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.ArtistDelete(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
