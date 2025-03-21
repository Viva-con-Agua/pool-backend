package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
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
	group.GET("/api_token", i.Get, vcago.KeyAuthMiddleware())
	group.GET("/:id", i.GetByID, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)
}

// Create
// @Security CookieAuth
// @Summary Create a Artist
// @Description creates an  Artist object.
// @Tags /events/artist
// @Accept json
// @Produce json
// @Param form body models.ArtistCreate true "Artist Data"
// @Model: vcago.Response
// @Success 201 {object} vcago.Response{payload=models.Artist}
// @Router /events/artist [post]
func (i *ArtistHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ArtistCreate)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Artist)
	if result, err = dao.ArtistInsert(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Created(result)
}

// Get
// @Security CookieAuth
// @Summary Get a List of Artists
// @Tags /events/artist
// @Accept json
// @Produce json
// @Param   q query   models.ArtistQuery   false  "string collection"  collectionFormat(multi)
// @Model: vcago.Response
// @Success 200 {object} vcago.Response{payload=[]models.Artist}
// @Failure 400 {object} vcago.Response{}
// @Router /events/artist [get]
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

// GetByID
// @Security CookieAuth
// @Summary Get a  Artist by ID
// @Tags /events/artist
// @Accept json
// @Produce json
// @Param id path string true "Artist ID"
// @Model: vcago.Response
// @Success 200 {object} vcago.Response{payload=models.Artist}
// @Router /events/artist/{id} [get]
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

// Update
// @Security CookieAuth
// @Summary Get a Artist by ID
// @Tags /events/artist
// @Accept json
// @Produce json
// @Param form body models.ArtistUpdate true "Artist Data"
// @Model: vcago.Response
// @Success 200 {object} vcago.Response{payload=models.Artist}
// @Router /events/artist [put]
func (i *ArtistHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ArtistUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Artist)
	if result, err = dao.ArtistUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Updated(result)
}

// DeleteByID
// @Security CookieAuth
// @Summary Get a  Artist by ID
// @Tags /events/artist
// @Accept json
// @Produce json
// @Param id path string true "Artist ID"
// @Success 200 {object} vmod.DeletedResponse
// @Router /events/artist/{id} [delete]
func (i *ArtistHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ArtistParam)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.ArtistDelete(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
