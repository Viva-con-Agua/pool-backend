package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

func CreateUserActive(c echo.Context) (err error) {
	ctx := c.Request().Context()
	user := new(vcapool.User)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	result := new(dao.UserActive)
	if result, err = result.Create(ctx, user); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("user_active", result).Created())
}
