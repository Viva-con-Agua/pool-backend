package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type DepositHandler struct {
	vcago.Handler
}

var Deposit = &DepositHandler{*vcago.NewHandler("deposit")}

func (i *DepositHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
}

func (i *DepositHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.DepositCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Deposit)
	if result, err = dao.DepositInsert(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Created(result)
}

func (i *DepositHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.DepositQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new([]models.Deposit)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if result, err = dao.DepositGet(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *DepositHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.DepositParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Deposit)
	if result, err = dao.DepositGetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *DepositHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.DepositUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Deposit)
	if result, err = dao.DepositUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Updated(result)
}

func (i *DepositHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.DepositParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	return c.Deleted(body.ID)

}
