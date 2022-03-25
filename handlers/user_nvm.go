package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func UserNVMConfirm(c echo.Context) (err error) {
	ctx := c.Request().Context()
	user := new(vcapool.AccessToken)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	if user.ActiveState != "confirmed" {
		return vcago.NewBadRequest("user_nvm", "active required")
	}
	if user.AddressID == "" {
		return vcago.NewBadRequest("user_nvm", "address required")
	}
	if user.Profile.Birthdate == 0 {
		return vcago.NewBadRequest("user_nvm", "birthdate required")
	}
	result := new(dao.UserNVM)
	if err = result.Get(ctx, bson.M{"user_id": user.ID}); err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	if vcago.MongoNoDocuments(err) {
		err = nil
		if result, err = result.Create(ctx, user.ID); err != nil {
			return
		}
	} else {
		if result, err = result.Confirm(ctx); err != nil {
			return
		}
	}
	return vcago.NewCreated("user_nvm", result)
}

func UserNVMReject(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.UserNVMRequest)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	//get requested user from token
	userReq := new(vcapool.AccessToken)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	//check if requested user has the network or operation permission
	if !userReq.PoolRoles.Validate("employee") {
		return vcago.NewPermissionDenied("user_nvm")
	}
	result := new(dao.UserNVM)
	if result, err = result.Reject(ctx, body.UserID); err != nil {
		return
	}
	return vcago.NewExecuted("user_nvm", result)
}

func UserNVMWithdraw(c echo.Context) (err error) {
	ctx := c.Request().Context()
	user := new(vcapool.AccessToken)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	result := new(dao.UserNVM)
	if err = result.Get(ctx, bson.M{"user_id": user.ID}); err != nil {
		return
	}
	if result, err = result.Withdraw(ctx); err != nil {
		return
	}
	return vcago.NewExecuted("user_nvm", result)
}
