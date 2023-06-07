package token

import (
	"log"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	vcago.Handler
}

var Role = &RoleHandler{*vcago.NewHandler("role")}

func (i *RoleHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.POST("/bulk", i.CreateBulk, accessCookie)
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
	result := new(vmod.Role)
	if result, err = dao.RoleInsert(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/crew/asp/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Created(result)
}

func (i *RoleHandler) CreateBulk(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.RoleBulkRequest)
	if c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(vmod.Role)
	if result, err = dao.RoleBulkUpdate(c.Ctx(), body, token); err != nil {
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
	result := new(vmod.Role)
	if result, err = dao.RoleDelete(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/crew/asp/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Deleted(body)
}
