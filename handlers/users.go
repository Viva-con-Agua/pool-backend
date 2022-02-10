package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

func CreateUser(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.User)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	if err = body.Create(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("users", body).Created())
}

func GetUser(c echo.Context) (err error) {
	ctx := c.Request().Context()
	result := new(dao.User)
	if err = result.Get(ctx, c.Param("id")); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("users", result).Selected())
}
