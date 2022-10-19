package admin

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	vcago.Handler
}

var Role = &RoleHandler{*vcago.NewHandler("role")}

func (i *RoleHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create)
	group.DELETE("", i.Delete)
}

func (i *RoleHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.RoleRequest)
	if c.BindAndValidate(body); err != nil {
		return
	}
	user := new(models.User)
	if err = dao.UserCollection.AggregateOne(
		c.Ctx(),
		models.UserPipeline().Match(body.MatchUser()).Pipe,
		user,
	); err != nil {
		return
	}
	var result *vmod.Role
	if result, err = body.New(); err != nil {
		return
	}
	if err = dao.PoolRoleCollection.InsertOne(c.Ctx(), result); err != nil {
		return
	}
	return c.Created(result)
}

func (i *RoleHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.RoleRequest)
	if c.BindAndValidate(body); err != nil {
		return
	}
	user := new(models.User)
	if err = dao.UserCollection.AggregateOne(
		c.Ctx(),
		models.UserPipeline().Match(body.MatchUser()).Pipe,
		user,
	); err != nil {
		return
	}
	result := new(vmod.Role)
	if err = dao.PoolRoleCollection.FindOne(
		c.Ctx(),
		body.Filter(),
		result,
	); err != nil {
		return
	}
	if err = dao.PoolRoleCollection.DeleteOne(c.Ctx(), body.Filter()); err != nil {
		return
	}
	return c.Deleted(body)
}
