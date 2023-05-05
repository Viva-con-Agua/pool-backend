package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
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
	result := body.Organizer()
	if err = dao.OrganizerCollection.InsertOne(c.Ctx(), result); err != nil {
		return c.ErrorResponse(err)
	}
	return c.Created(result)
}

func (i *OrganizerHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.OrganizerParam)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	result := new(models.Organizer)
	if err = dao.OrganizerCollection.FindOne(c.Ctx(), body.Filter(), result); err != nil {
		return c.ErrorResponse(err)
	}
	return c.Selected(result)
}

func (i *OrganizerHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.OrganizerUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	result := new(models.Organizer)
	if err = dao.OrganizerCollection.UpdateOne(c.Ctx(), body.Filter(), vmdb.UpdateSet(body), result); err != nil {
		return c.ErrorResponse(err)
	}
	return c.Updated(body)
}

func (i *OrganizerHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.OrganizerParam)
	if c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	if err = dao.OrganizerCollection.DeleteOne(c.Ctx(), body.Filter()); err != nil {
		return c.ErrorResponse(err)
	}
	return c.Deleted(body.ID)
}

func (i *OrganizerHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.OrganizerQuery)
	if err = c.BindAndValidate(body); err != nil {
		return c.ErrorResponse(err)
	}
	result := new([]models.Organizer)
	if err = dao.OrganizerCollection.Find(c.Ctx(), body.Filter(), result); err != nil {
		return c.ErrorResponse(err)
	}
	return c.Selected(result)
}
