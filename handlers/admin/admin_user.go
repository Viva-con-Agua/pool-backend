package admin

import (
	"pool-backend/dao"
	"pool-backend/models"

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
	group.POST("", i.Create)
	group.GET("/:id", i.GetByID)
	group.DELETE("/:id", i.Delete)
}

func (i *UserHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserDatabase)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new(models.User)
	if result, err = dao.UserInsert(c.Ctx(), body); err != nil {
		return
	}
	vcago.Nats.Publish("pool.user.created", result)
	return c.Created(result)
}

func (i *UserHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	if err = dao.UserDelete(c.Ctx(), body.ID); err != nil {
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
	if err = dao.UserCollection.AggregateOne(c.Ctx(), models.UserPipeline(false).Match(body.Match()).Pipe, result); err != nil {
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
	if err = dao.UserCollection.Aggregate(c.Ctx(), models.UserPipeline(false).Match(body.Filter()).Pipe, result); err != nil {
		return
	}
	return c.Selected(result)
}
