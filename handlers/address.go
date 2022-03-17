package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

func CreateAddress(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Address)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	user := new(vcapool.AccessToken)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	body.UserID = user.ID
	if err = body.Create(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("address", body).Created())
}

func GetAddress(c echo.Context) (err error) {
	ctx := c.Request().Context()
	result := new(dao.Address)
	if err = result.Get(ctx, c.Param("id")); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("address", result).Selected())
}

func UpdateAddress(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Address)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	if err = body.Update(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("address", body).Updated())
}

func DeleteAddress(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Address)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	if err = body.Delete(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("address", body).Deleted())
}

func ListAddress(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.AddressQuery)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	result := new(dao.AddressList)
	if err = result.Get(ctx, body.Filter()); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("address_list", result).Selected())
}
