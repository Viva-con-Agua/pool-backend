package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func AddressCreate(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Address)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userReq := new(vcapool.AccessToken)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	body.UserID = userReq.ID
	if err = body.Create(ctx); err != nil {
		return
	}
	return vcago.NewCreated("address", body)
}

func AddressGet(c echo.Context) (err error) {
	ctx := c.Request().Context()
	result := new(dao.Address)
	userReq := new(vcapool.AccessToken)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	var filter bson.M
	if userReq.Roles.Validate("employee;admin") {
		filter = bson.M{"_id": c.Param("id")}
	} else {
		filter = bson.M{"_id": c.Param("id"), "user_id": userReq.ID}
	}
	if err = result.Get(ctx, filter); err != nil {
		return
	}
	return vcago.NewSelected("address", result)
}

func AddressUpdate(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Address)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userReq := new(vcapool.AccessToken)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	if userReq.ID != body.UserID {
		return vcago.NewPermissionDenied("address", body.ID)
	}
	if err = body.Update(ctx); err != nil {
		return
	}
	return vcago.NewUpdated("address", body)
}

func AddressDelete(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Address)
	userReq := new(vcapool.AccessToken)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	if err = body.Delete(ctx, bson.M{"_id": c.Param("id"), "user_id": userReq.ID}); err != nil {
		return
	}
	return vcago.NewDeleted("address", c.Param("id"))
}

//TODO
func AddressList(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.AddressQuery)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	result := new(dao.AddressList)
	if err = result.Get(ctx, body.Filter()); err != nil {
		return
	}
	return vcago.NewSelected("address_list", result)
}
