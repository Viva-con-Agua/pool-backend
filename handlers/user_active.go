package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func RequestUserActive(c echo.Context) (err error) {
	ctx := c.Request().Context()
	user := new(vcapool.User)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}

	result := new(dao.UserActive)
	if err = result.Get(ctx, bson.M{"user_id": user.ID}); err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	if vcago.MongoNoDocuments(err) {
		err = nil
		if result, err = result.Create(ctx, user); err != nil {
			return
		}
	} else {
		if err = result.Request(ctx, user); err != nil {
			return
		}
	}
	return c.JSON(vcago.NewResponse("user_active", result).Created())
}

func ConfirmUserActive(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.UserActiveRequest)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	/*
		user := new(vcapool.User)
		if user, err = vcapool.AccessCookieUser(c); err != nil {
			return
		}*/
	result := new(dao.UserActive)
	if result, err = body.Confirm(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("user_active", result).Executed())
}

func RejectUserActive(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.UserActiveRequest)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	result := new(dao.UserActive)
	if result, err = body.Reject(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("user_active", result).Executed())
}
