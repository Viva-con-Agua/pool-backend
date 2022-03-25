package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type RoleRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func RoleCreate(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(RoleRequest)
	if vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userReq := new(vcapool.AccessToken)
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
	if user.NVM.Status != "confirmed" {
		return vcago.NewBadRequest("role", "nvm required", nil)
	}
	if !userReq.Roles.CheckRoot(role) && !userReq.PoolRoles.CheckRoot(role) {
		return vcago.NewBadRequest("role", "no permission for set this role", nil)
	}
	if err = (*dao.Role)(role).Create(ctx); err != nil {
		return
	}
	return vcago.NewCreated("role", role)
}

func RoleDelete(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(RoleRequest)
	if vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userReq := new(vcapool.AccessToken)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	user := new(dao.User)
	if err = user.Get(ctx, bson.M{"_id": body.UserID}); err != nil {
		return
	}
	role := new(dao.Role)
	if err = role.Get(ctx, bson.M{"name": body.Role, "user_id": body.UserID}); err != nil {
		return
	}
	if !userReq.Roles.CheckRoot((*vcago.Role)(role)) && !userReq.PoolRoles.CheckRoot((*vcago.Role)(role)) {
		return vcago.NewValidationError("no permission for delete this role")
	}
	if err = role.Delete(ctx); err != nil {
		return
	}
	return vcago.NewDeleted("role", role)
}
