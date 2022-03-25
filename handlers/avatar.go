package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func AvatarCreate(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Avatar)
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
	return vcago.NewCreated("avatar", body)
}

func AvatarDelete(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.Avatar)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	user := new(vcapool.AccessToken)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	if user.ID != body.UserID {
		return vcago.NewPermissionDenied("avatar", body.ID)
	}
	if err = body.Delete(ctx, bson.M{"_id": c.Param("id")}); err != nil {
		return
	}
	return vcago.NewDeleted("avatar", body)
}
