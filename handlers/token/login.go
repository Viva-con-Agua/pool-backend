package token

import (
	"net/http"
	"pool-user/dao"
	"pool-user/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type LoginHandler struct {
	vcago.Handler
}

var Login = &LoginHandler{*vcago.NewHandler("login")}

var HydraClient = vcago.NewHydraClient()

func (i *LoginHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("/callback", i.Callback)
	if vcago.Settings.Bool("API_TEST_LOGIN", "n", false) {
		group.POST("/testlogin", i.LoginAPI)
	}
	group.GET("/refresh", i.Refresh, refreshCookie)
	group.GET("/logout", i.Logout, accessCookie)
}

func (i *LoginHandler) Callback(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(vcago.Callback)
	if c.BindAndValidate(body); err != nil {
		return
	}
	tokenUser := new(vmod.User)
	if tokenUser, err = HydraClient.Callback(c.Ctx(), body); err != nil {
		return
	}
	result := new(models.User)
	if err = dao.UserCollection.AggregateOne(
		c.Ctx(),
		models.UserPipeline().Match(models.UserMatch(tokenUser.ID)).Pipe,
		result,
	); err != nil && vmdb.ErrNoDocuments(err) {
		return
	}
	if vmdb.ErrNoDocuments(err) {
		userDatabase := models.NewUserDatabase(tokenUser)
		if err = dao.UserCollection.InsertOne(c.Ctx(), userDatabase); err != nil {
			return
		}
		if err = dao.UserCollection.AggregateOne(
			c.Ctx(),
			models.UserPipeline().Match(models.UserMatch(tokenUser.ID)).Pipe,
			result,
		); err != nil {
			return
		}
		vcago.Nats.Publish("user.created", result)
	}
	if tokenUser.CheckUpdate(result.LastUpdate) {
		userUpdate := models.NewUserUpdate(tokenUser)
		if err = dao.UserCollection.UpdateOne(c.Ctx(), userUpdate.Filter(), vmdb.NewUpdateSet(userUpdate), result); err != nil {
			return
		}
		vcago.Nats.Publish("user.updated", result)
	}
	token := new(vcago.AuthToken)
	if token, err = result.AuthToken(); err != nil {
		return
	}

	c.SetCookie(token.AccessCookie())
	c.SetCookie(token.RefreshCookie())
	return c.Selected(result)
}

func (i *LoginHandler) LoginAPI(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserEmail)
	if c.BindAndValidate(body); err != nil {
		return
	}
	result := new(models.User)
	if err = dao.UserCollection.AggregateOne(
		c.Ctx(),
		models.UserPipeline().Match(models.UserMatchEmail(body.Email)).Pipe,
		result,
	); err != nil && vmdb.ErrNoDocuments(err) {
		return
	}
	token := new(vcago.AuthToken)
	if token, err = result.AuthToken(); err != nil {
		return
	}
	c.SetCookie(token.AccessCookie())
	c.SetCookie(token.RefreshCookie())
	return c.Selected(result)
}

func (i *LoginHandler) Refresh(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	var userID string
	if userID, err = c.RefreshTokenID(); err != nil {
		return
	}
	result := new(models.User)
	if err = dao.UserCollection.AggregateOne(
		c.Ctx(),
		models.UserPipeline().Match(models.UserMatch(userID)).Pipe,
		result,
	); err != nil && vmdb.ErrNoDocuments(err) {
		return
	}
	token := new(vcago.AuthToken)
	if token, err = result.AuthToken(); err != nil {
		return
	}
	c.SetCookie(token.AccessCookie())
	c.SetCookie(token.RefreshCookie())
	return c.Selected(result)
}

func (i *LoginHandler) Logout(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	c.SetCookie(vcago.ResetAccessCookie())
	c.SetCookie(vcago.ResetRefreshCookie())
	return c.SuccessResponse(http.StatusOK, "logout", "user", nil)
}
