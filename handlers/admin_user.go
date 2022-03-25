package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func UserCreateAdmin(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.UserInsert)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	if err = body.Create(ctx); err != nil {
		return
	}
	return vcago.NewCreated("users", body)
}

func UserGetAdmin(c echo.Context) (err error) {
	ctx := c.Request().Context()
	result := new(dao.User)
	if err = result.Get(ctx, bson.M{"_id": c.Param("id")}); err != nil {
		return
	}
	return vcago.NewSelected("users", result)
}

func UserListAdmin(c echo.Context) (err error) {
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
