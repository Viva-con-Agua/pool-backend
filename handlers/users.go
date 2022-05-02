package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type UserHandler struct {
	vcago.Handler
}

func NewUserHandler() *UserHandler {
	handler := vcago.NewHandler("user")
	return &UserHandler{
		*handler,
	}
}

func (UserHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.UserDelete(c.Ctx(), token); err != nil {
		return
	}
	return c.Deleted(token.ID)
}

func UserCreate(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.UserDatabase)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	result := new(vcapool.User)
	if result, err = body.Create(ctx); err != nil {
		return
	}
	vcago.Nats.Publish("user.created", result)
	return vcago.NewCreated("users", body)
}

func UserGet(c echo.Context) (err error) {
	ctx := c.Request().Context()
	result := new(dao.User)
	if err = result.Get(ctx, bson.M{"_id": c.Param("id")}); err != nil {
		return
	}
	return vcago.NewSelected("users", result)
}

func UserList(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.UserQuery)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	result := new(dao.UserList)
	if result, err = body.List(ctx); err != nil {
		return
	}
	return vcago.NewSelected("user_list", result)
}
