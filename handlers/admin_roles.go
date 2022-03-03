package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type AdminRoleRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func AdminRoleCreate(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(AdminRoleRequest)
	if vcago.BindAndValidate(c, body); err != nil {
		return
	}
	user := new(dao.User)
	if err = user.Get(ctx, bson.M{"email": body.Email}); err != nil {
		return
	}
	var role *vcago.Role
	if role, err = vcapool.NewRole(body.Role, user.ID); err != nil {
		return
	}
	if err = (*dao.Role)(role).Create(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("role", role).Created())
}

func AdminRoleDelete(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(AdminRoleRequest)
	if vcago.BindAndValidate(c, body); err != nil {
		return
	}
	user := new(dao.User)
	if err = user.Get(ctx, bson.M{"email": body.Email}); err != nil {
		return
	}
	role := new(dao.Role)
	if err = role.Get(ctx, bson.M{"name": body.Role}); err != nil {
		return
	}
	if err = role.Delete(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("role", role).Deleted())
}
