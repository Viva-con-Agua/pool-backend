package handlers

import (
	"errors"
	"pool-user/dao"
	"time"

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
	if time.Unix(user.Modified.Created, 0).AddDate(0, 6, 0).Unix() < time.Now().Unix() {
		return vcago.NewStatusBadRequest(errors.New("create_date"))
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

//ConfirmUserActive is the webapp handler for confirm the active state of an user.
func ConfirmUserActive(c echo.Context) (err error) {
	ctx := c.Request().Context()
	//validate and bind body
	body := new(dao.UserActiveRequest)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	//get requested user from token
	userReq := new(vcapool.User)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	//check if requested user has the network or operation permission
	if !userReq.PoolRoles.Validate("network;operation") {
		return vcago.NewStatusPermissionDenied()
	}
	//check if requested user and the target users has the same crew
	userCrew := new(dao.UserCrew)
	if err = userCrew.Permission(ctx, bson.M{"user_id": body.UserID, "crew_id": userReq.Crew}); err != nil {
		return
	}
	//confirm active state
	result := new(dao.UserActive)
	if result, err = body.Confirm(ctx); err != nil {
		return
	}
	//response the result as vcago.Response
	return c.JSON(vcago.NewResponse("user_active", result).Executed())
}

func RejectUserActive(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.UserActiveRequest)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	//get requested user from token
	userReq := new(vcapool.User)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	//check if requested user has the network or operation permission
	if !userReq.PoolRoles.Validate("network;operation") {
		return vcago.NewStatusPermissionDenied()
	}
	//check if requested user and the target users has the same crew
	userCrew := new(dao.UserCrew)
	if err = userCrew.Permission(ctx, bson.M{"user_id": body.UserID, "crew_id": userReq.Crew}); err != nil {
		return
	}
	result := new(dao.UserActive)
	if result, err = body.Reject(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("user_active", result).Executed())
}

func WithdrawUserActive(c echo.Context) (err error) {
	ctx := c.Request().Context()
	user := new(vcapool.User)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	result := new(dao.UserActive)
	if err = result.Get(ctx, bson.M{"user_id": user.ID}); err != nil {
		return
	}
	if result, err = result.Withdraw(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("user_active", result).Executed())
}
