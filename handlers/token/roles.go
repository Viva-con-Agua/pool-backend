package token

import (
	"log"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type RoleHandler struct {
	vcago.Handler
}

var Role = &RoleHandler{*vcago.NewHandler("role")}

func (i *RoleHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.DELETE("", i.Delete, accessCookie)
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
	if err = dao.UserCollection.AggregateOne(
		c.Ctx(),
		models.UserPipeline(false).Match(body.MatchUser()).Pipe,
		user,
	); err != nil {
		log.Print(err)
		return
	}
	var result *vmod.Role
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
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/crew/asp/"); err != nil {
			log.Print(err)
		}
	}()
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
		bson.D{{Key: "_id", Value: body.UserID}},
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
	if !token.Roles.CheckRoot((*vmod.Role)(result)) && !token.PoolRoles.CheckRoot((*vmod.Role)(result)) {
		return vcago.NewValidationError("no permission for delete this role")
	}
	if err = dao.PoolRoleCollection.DeleteOne(c.Ctx(), body.Filter()); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/crew/asp/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Deleted(body)
}
