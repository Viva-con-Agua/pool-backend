package token

import (
	"net/http"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type NVMHandler struct {
	vcago.Handler
}

var NVM = &NVMHandler{*vcago.NewHandler("nvm")}

func (i *NVMHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.GET("/confirm", i.Confirm, accessCookie)
	group.POST("/confirm/:id", i.ConfirmUser, accessCookie)
	group.DELETE("/reject/:id", i.RejectUser, accessCookie)
	group.GET("/withdraw", i.Withdraw, accessCookie)
}

func (i *NVMHandler) Confirm(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	var result *models.NVM
	if result, err = dao.NVMConfirm(c.Ctx(), token); err != nil {
		return
	}
	return c.SuccessResponse(http.StatusOK, "successfully_confirmed", "nvm", result)
}

func (i *NVMHandler) ConfirmUser(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.NVMIDParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	//get requested user from token
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	var result *models.NVM
	if result, err = dao.NVMConfirmUser(c.Ctx(), body, token); err != nil {
		return
	}
	return c.SuccessResponse(http.StatusOK, "successfully_confirmed", "nvm", result)
}

func (i *NVMHandler) RejectUser(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.NVMIDParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	var result *models.NVM
	if result, err = dao.NVMRejectUser(c.Ctx(), body, token); err != nil {
		return
	}
	return c.SuccessResponse(http.StatusOK, "successfully_rejected", "nvm", result)
}

func (i *NVMHandler) Withdraw(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	var result *models.NVM
	if result, err = dao.NVMWithdraw(c.Ctx(), token); err != nil {
		return
	}
	return c.SuccessResponse(http.StatusOK, "successfully_withdrawed", "nvm", result)
}
