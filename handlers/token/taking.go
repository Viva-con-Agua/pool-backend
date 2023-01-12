package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type TakingHandler struct {
	vcago.Handler
}

var Taking = &TakingHandler{*vcago.NewHandler("taking")}

func (i *TakingHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)

}

func (i *TakingHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.TakingCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	var result *models.Taking
	if result, err = dao.TakingInsert(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Created(result)
}

func (i TakingHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	body := new(models.TakingUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new(models.Taking)
	if result, err = dao.TakingUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Updated(result)
}

func (i TakingHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.TakingQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	var result *[]models.Taking
	if result, err = dao.TakingGet(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Listed(result)
}

func (i TakingHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(vmod.IDParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	var result *models.Taking
	if result, err = dao.TakingGetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i TakingHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(vmod.IDParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	var result *vmod.DeletedResponse
	if result, err = dao.TakingDeletetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Deleted(result)
}
