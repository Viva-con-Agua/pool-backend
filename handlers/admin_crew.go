package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

func ListCrewAdmin(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.CrewQuery)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	result := new(dao.CrewList)
	if result, err = body.List(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("crew_list", result).Selected())
}
