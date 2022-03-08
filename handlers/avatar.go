package handlers

import (
	"errors"
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

func CreateAvatar(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Profile)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	user := new(vcapool.User)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	body.UserID = user.ID
	if err = body.Create(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("avatar", body).Created())
}

func DeleteAvatar(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Profile)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	user := new(vcapool.User)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	if user.ID != body.UserID {
		return vcago.NewStatusBadRequest(errors.New("permission denied"))
	}
	if err = body.Delete(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("avatar", body).Deleted())

}
