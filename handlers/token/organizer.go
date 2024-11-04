package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type OrganizerHandler struct {
	vcago.Handler
}

var Organizer = &OrganizerHandler{*vcago.NewHandler("organizer")}

func (i *OrganizerHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)
}

func (i *OrganizerHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.OrganizerCreate)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Organizer)
	if result, err = dao.OrganizerInsert(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Created(result)
}

func (i *OrganizerHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.OrganizerQuery)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	result := new([]models.Organizer)
	if result, err = dao.OrganizerGet(c.Ctx(), body); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *OrganizerHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.OrganizerParam)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	result := new(models.Organizer)
	if result, err = dao.OrganizerGetByID(c.Ctx(), body); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *OrganizerHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.OrganizerUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Organizer)
	if result, err = dao.OrganizerUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Updated(result)
}

func (i *OrganizerHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.OrganizerParam)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.OrganizerDelete(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
