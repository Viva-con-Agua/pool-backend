package token

import (
	"pool-user/dao"
	"pool-user/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	vcago.Handler
}

var Role = &RoleHandler{*vcago.NewHandler("role")}

func (i *RoleHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, vcapool.AccessCookieConfig())
	group.DELETE("", i.Delete, vcapool.AccessCookieConfig())
}

func (i *RoleHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.RoleRequest)
	if c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	user := new(models.User)
	if err = dao.UserCollection.FindOne(
		c.Ctx(),
		models.UserPipeline().Match(body.MatchUser()).Pipe,
		user,
	); err != nil {
		return
	}
	var result *vcago.Role
	if result, err = body.New(); err != nil {
		return
	}
	if user.NVM.Status != "confirmed" {
		return vcago.NewBadRequest("role", "nvm required", nil)
	}
	if !token.Roles.CheckRoot(result) && !token.PoolRoles.CheckRoot(result) {
		return vcago.NewBadRequest("role", "no permission for set this role", nil)
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
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	user := new(models.User)
	if err = dao.UserCollection.FindOne(
		c.Ctx(),
		models.UserPipeline().Match(body.Match()).Pipe,
		user,
	); err != nil {
		return
	}
	result := new(vcago.Role)
	if err = dao.PoolRoleCollection.FindOne(
		c.Ctx(),
		vmdb.NewPipeline().Match(body.Match()).Pipe,
		result,
	); err != nil {
		return
	}
	if !token.Roles.CheckRoot((*vcago.Role)(result)) && !token.PoolRoles.CheckRoot((*vcago.Role)(result)) {
		return vcago.NewValidationError("no permission for delete this role")
	}
	if err = dao.PoolRoleCollection.DeleteOne(c.Ctx(), body.Filter()); err != nil {
		return
	}
	return c.Deleted(body)
}
