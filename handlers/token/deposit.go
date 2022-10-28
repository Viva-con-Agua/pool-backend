package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type DepositHandler struct {
	vcago.Handler
}

var Deposit = &DepositHandler{*vcago.NewHandler("deposit")}

func (i *DepositHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.GET("", i.Get, accessCookie)
}

func (i *DepositHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.DepositCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	var result *models.Deposit
	if result, err = dao.DepositInsert(c.Ctx(), body); err != nil {
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
	var result *[]models.Deposit
	if result, err = dao.DepositGet(c.Ctx(), body); err != nil {
		return
	}
	return c.Selected(result)
}
