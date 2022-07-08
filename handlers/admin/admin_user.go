package admin

import (
	"pool-user/dao"
	"pool-user/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	vcago.Handler
}

var User = &UserHandler{*vcago.NewHandler("user")}

func (i *UserHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.GET("", i.Get)
	group.GET("/:id", i.GetByID)
	group.DELETE("/:id", i.Delete)
}

/*
func (i *UserHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(dao.UserDatabase)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new(models.User)
	return vcago.NewCreated("user", result)
}*/

func (i *UserHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	if err = dao.UserCollection.DeleteOne(c.Ctx(), body.FilterAdmin()); err != nil {
		return
	}
	return c.Deleted(body.ID)
}

func (i *UserHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new(models.User)
	if err = dao.UserCollection.AggregateOne(c.Ctx(), body.Pipeline(), result); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *UserHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new([]models.User)
	if err = dao.UserCollection.Aggregate(c.Ctx(), body.Pipeline(), result); err != nil {
		return
	}
	return c.Selected(result)
}
