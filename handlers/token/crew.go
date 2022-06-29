package token

import (
	"pool-user/dao"
	"pool-user/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type CrewHandler struct {
	vcago.Handler
}

var Crew = &CrewHandler{*vcago.NewHandler("crew")}

func (i *CrewHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, vcapool.AccessCookieConfig())
	group.GET("", i.Get)
	group.GET("/:id", i.GetByID)
	group.PUT("", i.Update, vcapool.AccessCookieConfig())
	group.DELETE("/:id", i.Delete, vcapool.AccessCookieConfig())
}

func (i *CrewHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = models.CrewPermission(token); err != nil {
		return
	}
	result := body.Crew()
	if err = dao.CrewsCollection.InsertOne(c.Ctx(), result); err != nil {
		c.Log(err)(err)
		return c.ErrorResponse(err)
	}
	return c.Created(result)
}

func (i *CrewHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new(models.Crew)
	if dao.CrewsCollection.Find(c.Ctx(), body.Pipeline(), result); err != nil {
		return
	}
	return c.Listed(result)
}

func (i *CrewHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new(models.Crew)
	if err = dao.CrewsCollection.FindOne(c.Ctx(), body.Pipeline(), result); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *CrewHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = models.CrewPermission(token); err != nil {
		return
	}
	result := new(models.Crew)
	if err = dao.CrewsCollection.UpdateOne(c.Ctx(), body.Filter(), vmdb.NewUpdateSet(body), result); err != nil {
		return
	}
	return vcago.NewUpdated("crew", body)
}

func (i *CrewHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = models.CrewPermission(token); err != nil {
		return
	}
	if err = dao.CrewsCollection.DeleteOne(c.Ctx(), body.Filter()); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
