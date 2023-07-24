package token

import (
	"net/http"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

// ActiveHandler represents the vcago.Handler for the Active model.
type ActiveHandler struct {
	vcago.Handler
}

// Active is an ActiveHandler.
var Active = &ActiveHandler{*vcago.NewHandler("active")}

func (i *ActiveHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.GET("/request", i.Request, accessCookie)
	group.POST("/confirm", i.Confirm, accessCookie)
	group.POST("/reject", i.Reject, accessCookie)
	group.GET("/withdraw", i.Withdraw, accessCookie)
}

// Request handles an active request call.
func (i *ActiveHandler) Request(cc echo.Context) (err error) {
	//load context
	c := cc.(vcago.Context)
	//get access token
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Active)
	if result, err = dao.ActiveRequest(c.Ctx(), token); err != nil {
		return
	}
	//return "successfully requested."
	return c.SuccessResponse(http.StatusOK, "successfully_requested", "active", result)
}

// Confirm handles an active confirm request.
func (i *ActiveHandler) Confirm(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	//validate and bind body
	body := new(models.ActiveParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	//get requested user from token
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Active)
	if result, err = dao.ActiveConfirm(c.Ctx(), body, token); err != nil {
		return
	}
	dao.ActiveNotification(c.Ctx(), result, "active_confirm")
	return c.SuccessResponse(http.StatusOK, "successfully_confirmed", "active", result)
}

func (i *ActiveHandler) Reject(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ActiveParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	//get requested user from token
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Active)
	if result, err = dao.ActiveReject(c.Ctx(), body, token); err != nil {
		return
	}
	dao.ActiveNotification(c.Ctx(), result, "active_reject")
	return c.SuccessResponse(http.StatusOK, "successfully_rejected", "active", result)
}

func (i *ActiveHandler) Withdraw(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	//reject nvm state
	result := new(models.Active)
	if result, err = dao.ActiveWithdraw(c.Ctx(), token); err != nil {
		return
	}
	return c.SuccessResponse(http.StatusOK, "successfully_withdrawn", "active", result)
}
