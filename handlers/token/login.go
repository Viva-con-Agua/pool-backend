package token

import (
	"log"
	"net/http"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
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
	if dao.TestLogin {
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
		return vcago.NewBadRequest("Error in callback. Maybe testlogin and oidc skip is enabeled in .env?", err.Error())
	}
	result := new(models.User)
	if err = dao.UserCollection.AggregateOne(
		c.Ctx(),
		models.UserPipeline(true).Match(models.UserMatch(tokenUser.ID)).Pipe,
		result,
	); err != nil && !vmdb.ErrNoDocuments(err) {
		return
	}
	if vmdb.ErrNoDocuments(err) {

		err = nil
		userDatabase := models.NewUserDatabase(tokenUser)
		if result, err = dao.UserInsert(c.Ctx(), userDatabase); err != nil {
			return
		}
		go func() {
			if err = dao.IDjango.Post(result, "/v1/pool/user/"); err != nil {
				log.Print(err)
			}
		}()
		vcago.Nats.Publish("pool.user.created", result)
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
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new(models.User)
	if err = dao.UserCollection.AggregateOne(
		c.Ctx(),
		models.UserPipeline(true).Match(models.UserMatchEmail(body.Email)).Pipe,
		result,
	); err != nil {
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
		models.UserPipeline(true).Match(models.UserMatch(userID)).Pipe,
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
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	c.SetCookie(vcago.ResetAccessCookie())
	c.SetCookie(vcago.ResetRefreshCookie())
	return c.SuccessResponse(http.StatusOK, "logout", "user", nil)
}
