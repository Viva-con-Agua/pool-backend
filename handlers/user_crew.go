package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

func CreateUserCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.UserCrewCreateRequest)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	user := new(vcapool.User)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	result := new(dao.UserCrew)
	if result, err = body.Create(ctx, user.ID); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("user_crew", result).Created())
}
