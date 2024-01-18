package token

import (
	"log"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type RoleHistoryHandler struct {
	vcago.Handler
}

var RoleHistory = &RoleHistoryHandler{*vcago.NewHandler("role_history")}

func (i *RoleHistoryHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.POST("/bulk", i.CreateBulk, accessCookie)
	group.POST("/confirm", i.ConfirmSelection, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.DELETE("", i.Delete, accessCookie)
}

func (i *RoleHistoryHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.RoleHistoryCreate)
	if c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.RoleHistory)
	if result, err = dao.RoleHistoryInsert(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Created(result)
}

func (i *RoleHistoryHandler) CreateBulk(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.RoleHistoryBulkRequest)
	if c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.RoleBulkExport)
	if result, err = dao.RoleHistoryBulkInsert(c.Ctx(), body, token); err != nil {
		return
	}
	if _, err = dao.CrewUpdateAspSelection(c.Ctx(), &models.CrewParam{ID: body.CrewID}, "selected", token); err != nil {
		return
	}
	dao.RoleHistoryAdminNotification(c.Ctx(), &models.CrewParam{ID: body.CrewID})
	return c.Created(result)
}

func (i *RoleHistoryHandler) ConfirmSelection(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.RoleHistoryRequest)
	if c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	body.Confirmed = true
	history := new([]models.RoleHistory)
	if history, err = dao.RoleHistoryConfirm(c.Ctx(), body, token); err != nil {
		return
	}
	result := new(models.RoleBulkExport)
	userRolesMap := make(map[string]*models.BulkUserRoles)
	if result, userRolesMap, err = dao.RoleBulkConfirm(c.Ctx(), history, body.CrewID, token); err != nil {
		return
	}
	if _, err = dao.CrewUpdateAspSelection(c.Ctx(), &models.CrewParam{ID: body.CrewID}, "inactive", token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/asps/"); err != nil {
			log.Print(err)
		}
	}()
	if err = dao.RoleNotification(c.Ctx(), userRolesMap); err != nil {
		return
	}
	return c.Created(result)
}

func (i *RoleHistoryHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.RoleHistoryRequest)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new([]models.RoleHistory)
	var listSize int64
	if result, listSize, err = dao.RoleHistoryGet(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Listed(result, listSize)
}

func (i *RoleHistoryHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.RoleHistoryRequest)
	if c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.RoleHistory)
	if result, err = dao.RoleHistoryDelete(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/crew/asp/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Deleted(body)
}
