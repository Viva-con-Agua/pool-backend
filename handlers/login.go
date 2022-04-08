package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type LoginHandler struct {
	vcago.Handler
}

func NewLoginHandler() *LoginHandler {
	handler := vcago.NewHandler("login")
	return &LoginHandler{
		*handler,
	}
}

var HydraClient = vcago.NewHydraClient()

func (i *LoginHandler) Callback(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(vcago.Callback)
	if vcago.BindAndValidate(c, body); err != nil {
		return
	}
	tokenUser := new(vcago.User)
	if tokenUser, err = HydraClient.Callback(c.Ctx(), body); err != nil {
		return
	}
	userSystem := dao.NewUserSystem(tokenUser)
	result := new(vcapool.User)
	if result, err = userSystem.Get(c.Ctx()); err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	if vcago.MongoNoDocuments(err) {
		if result, err = userSystem.Create(c.Ctx()); err != nil {
			return
		}
	}
	if tokenUser.CheckUpdate(result.LastUpdate) {
		if result, err = userSystem.Update(c.Ctx()); err != nil {
			return
		}
	}
	token := new(vcapool.AuthToken)
	if token, err = vcapool.NewAuthToken(result); err != nil {
		return
	}
	c.SetCookie(token.AccessCookie())
	c.SetCookie(token.RefreshCookie())
	return c.Selected(result)
}

type LoginAPIRequest struct {
	Email string `json:"email"`
}

func (i *LoginHandler) LoginAPI(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(LoginAPIRequest)
	if vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userSystem := new(dao.User)
	result := new(vcapool.User)
	if err = userSystem.Get(c.Ctx(), bson.M{"email": body.Email}); err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	result = (*vcapool.User)(userSystem)
	token := new(vcapool.AuthToken)
	if token, err = vcapool.NewAuthToken(result); err != nil {
		return
	}
	c.SetCookie(token.AccessCookie())
	c.SetCookie(token.RefreshCookie())
	return c.Selected(result)
}

func (i *LoginHandler) Refresh(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	var userID string
	if userID, err = vcapool.RefreshCookieUserID(c); err != nil {
		return
	}
	userDAO := new(dao.User)
	if err = userDAO.Get(c.Ctx(), bson.M{"_id": userID}); err != nil {
		return
	}
	user := vcapool.User(*userDAO)
	token := new(vcapool.AuthToken)
	if token, err = vcapool.NewAuthToken(&user); err != nil {
		return
	}
	c.SetCookie(token.AccessCookie())
	c.SetCookie(token.RefreshCookie())
	return vcago.NewSelected("refresh_token", user)
}

func (i *LoginHandler) Logout(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	c.SetCookie(vcapool.ResetAccessCookie())
	c.SetCookie(vcapool.ResetRefreshCookie())
	return vcago.NewExecuted("logout", token.ID)
}
