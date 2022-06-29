package token

import (
	"net/http"
	"pool-user/dao"
	"pool-user/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
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
	group.GET("/refresh", i.Refresh, vcapool.RefreshCookieConfig())
	group.GET("/logout", i.Logout, vcago.AccessCookieMiddleware(&vcapool.AccessToken{}))
}

func (i *LoginHandler) Callback(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(vcago.Callback)
	if c.BindAndValidate(body); err != nil {
		return
	}
	tokenUser := new(vcago.User)
	if tokenUser, err = HydraClient.Callback(c.Ctx(), body); err != nil {
		return
	}
	result := new(models.User)
	if err = dao.UserCollection.FindOne(
		c.Ctx(),
		models.UserPipeline().Match(models.UserMatch(tokenUser.ID)).Pipe,
		result,
	); err != nil && err != mongo.ErrNoDocuments {
		return
	}
	if err == mongo.ErrNoDocuments {
		err = nil
		userDatabase := models.NewUserDatabase(tokenUser)
		if err = dao.UserCollection.InsertOne(c.Ctx(), userDatabase); err != nil {
			return
		}
		if err = dao.UserCollection.FindOne(
			c.Ctx(),
			models.UserPipeline().Match(models.UserMatch(tokenUser.ID)).Pipe,
			result,
		); err != nil {
			return
		}
	}
	if tokenUser.CheckUpdate(result.LastUpdate) {
		userUpdate := models.NewUserUpdate(tokenUser)
		if err = dao.UserCollection.UpdateOne(c.Ctx(), userUpdate.Filter(), vmdb.NewUpdateSet(userUpdate), result); err != nil {
			return
		}
	}
	token := new(vcapool.AuthToken)
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
	if err = dao.UserCollection.FindOne(
		c.Ctx(),
		models.UserPipeline().Match(models.UserMatchEmail(body.Email)).Pipe,
		result,
	); err != nil && err != mongo.ErrNoDocuments {
		return
	}
	token := new(vcapool.AuthToken)
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
	if userID, err = vcapool.RefreshCookieUserID(c); err != nil {
		return
	}
	result := new(models.User)
	if err = dao.UserCollection.FindOne(
		c.Ctx(),
		models.UserPipeline().Match(models.UserMatch(userID)).Pipe,
		result,
	); err != nil && err != mongo.ErrNoDocuments {
		return
	}
	token := new(vcapool.AuthToken)
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
	c.SetCookie(vcapool.ResetAccessCookie())
	c.SetCookie(vcapool.ResetRefreshCookie())
	return c.SuccessResponse(http.StatusOK, "logout", "user", nil)
}
