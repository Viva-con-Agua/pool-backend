package admin

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type CrewHandler struct {
	vcago.Handler
}

var Crew = &CrewHandler{*vcago.NewHandler("crew")}

func (i *CrewHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create)
	group.GET("", i.Get)
	group.DELETE("/:id", i.Delete)
}

func (i *CrewHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new(models.Crew)
	if result, err = dao.CrewInsert(c.Ctx(), body); err != nil {
		return
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
	if dao.CrewsCollection.Find(c.Ctx(), body.Filter(), result); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *CrewHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	if err = dao.CrewDelete(c.Ctx(), body); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
