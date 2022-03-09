package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

func CreateCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Crew)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userReq := new(vcapool.User)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	if !userReq.Roles.Validate("employee") {
		return vcago.NewStatusPermissionDenied()
	}
	if err = body.Create(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("crew", body).Created())
}

func GetCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	result := new(dao.Crew)
	if err = result.Get(ctx, c.Param("id")); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("crew", result).Selected())
}

func UpdateCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Crew)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userReq := new(vcapool.User)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	if !userReq.Roles.Validate("employee") {
		return vcago.NewStatusPermissionDenied()
	}
	if err = body.Update(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("crew", body).Updated())
}

func DeleteCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Crew)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userReq := new(vcapool.User)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	if !userReq.Roles.Validate("employee") {
		return vcago.NewStatusPermissionDenied()
	}
	if err = body.Delete(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("crew", body).Deleted())
}

func ListCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.CrewQuery)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	result := new(dao.CrewList)
	if err = result.Get(ctx, body.Filter()); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("crew_list", result).Selected())
}
