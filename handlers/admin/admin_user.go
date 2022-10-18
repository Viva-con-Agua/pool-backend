package admin

import (
	"pool-user/dao"
	"pool-user/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
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
	if err = dao.UserCollection.InsertOne(c.Ctx(), body); err != nil {
		return
	}
	result := new(models.User)
	if err = dao.UserCollection.FindOne(
		c.Ctx(),
		bson.D{{Key: "_id", Value: body.ID}},
		result,
	); err != nil {
		return
	}
	vcago.Nats.Publish("pool.user.creatd", result)
	return c.Created(result)
}

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
