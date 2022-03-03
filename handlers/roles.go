package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type RoleRequest struct {
	UserID string `json:"id"`
	Role   string `json:"role"`
}

func RoleCreate(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(RoleRequest)
	if vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userReq := new(vcapool.User)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	user := new(dao.User)
	if err = user.Get(ctx, bson.M{"_id": body.UserID}); err != nil {
		return
	}
	var role *vcago.Role
	if role, err = vcapool.NewRole(body.Role, user.ID); err != nil {
		return
	}
	if !userReq.Roles.CheckRoot(role) && !userReq.PoolRoles.CheckRoot(role) {
		return vcago.NewValidationError("no permission for set this role")
	}
	if err = (*dao.Role)(role).Create(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("role", role).Created())
}

func RoleDelete(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(RoleRequest)
	if vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userReq := new(vcapool.User)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	user := new(dao.User)
	if err = user.Get(ctx, bson.M{"_id": body.UserID}); err != nil {
		return
	}
	var role *dao.Role
	if err = role.Get(ctx, bson.M{"name": body.Role}); err != nil {
		return
	}
	if !userReq.Roles.CheckRoot((*vcago.Role)(role)) && !userReq.PoolRoles.CheckRoot((*vcago.Role)(role)) {
		return vcago.NewValidationError("no permission for set this role")
	}
	if err = role.Delete(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("role", role).Deleted())
}
