package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

func CrewCreate(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Crew)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userReq := new(vcapool.AccessToken)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	if !userReq.Roles.Validate("employee") {
		return vcago.NewPermissionDenied("crew", nil)
	}
	if err = body.Create(ctx); err != nil {
		return
	}
	return vcago.NewCreated("crew", body)
}

func CrewGet(c echo.Context) (err error) {
	ctx := c.Request().Context()
	result := new(dao.Crew)
	if err = result.Get(ctx, c.Param("id")); err != nil {
		return
	}
	return vcago.NewSelected("crew", result)
}

func CrewUpdate(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Crew)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userReq := new(vcapool.AccessToken)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	if !userReq.Roles.Validate("employee") {
		return vcago.NewPermissionDenied("crew")
	}
	if err = body.Update(ctx); err != nil {
		return
	}
	return vcago.NewUpdated("crew", body)
}

func CrewDelete(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Crew)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userReq := new(vcapool.AccessToken)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	if !userReq.Roles.Validate("employee") {
		return vcago.NewPermissionDenied("crew")
	}
	if err = body.Delete(ctx); err != nil {
		return
	}
	return vcago.NewDeleted("crew", body)
}

func CrewList(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.CrewQuery)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	result := new(dao.CrewList)
	if result, err = body.List(ctx); err != nil {
		return
	}
	return vcago.NewSelected("crew_list", result)
}
