package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

var HydraClient = vcago.NewHydraClient()

func CallbackHandler(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(vcago.Callback)
	if vcago.BindAndValidate(c, body); err != nil {
		return
	}
	tokenUser := new(vcago.User)
	if tokenUser, err = HydraClient.Callback(ctx, body); err != nil {
		return
	}
	userDAO := new(dao.User)
	if err = userDAO.Get(ctx, bson.M{"_id": tokenUser.ID}); err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	if vcago.MongoNoDocuments(err) {
		userInsert := dao.NewUser(tokenUser)
		if err = userInsert.Create(ctx); err != nil {
			return
		}
		if err = userDAO.Get(ctx, bson.M{"_id": tokenUser.ID}); err != nil && !vcago.MongoNoDocuments(err) {
			return
		}
	}
	if tokenUser.CheckUpdate(userDAO.LastUpdate) {
		userInsert := dao.ConvertUser(tokenUser, &userDAO.Modified)
		if err = userInsert.Update(ctx); err != nil {
			return
		}
		if err = userDAO.Get(ctx, bson.M{"_id": tokenUser.ID}); err != nil {
			return
		}
	}
	user := vcapool.User(*userDAO)
	token := new(vcapool.AuthToken)
	if token, err = vcapool.NewAuthToken(&user); err != nil {
		return vcago.NewStatusInternal(err)
	}
	c.SetCookie(token.AccessCookie())
	c.SetCookie(token.RefreshCookie())
	return c.JSON(vcago.NewResponse("access_user", user).Selected())
}

func RefreshHandler(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var userID string
	if userID, err = vcapool.RefreshCookieUserID(c); err != nil {
		return
	}
	userDAO := new(dao.User)
	if err = userDAO.Get(ctx, bson.M{"_id": userID}); err != nil {
		return
	}
	user := vcapool.User(*userDAO)
	token := new(vcapool.AuthToken)
	if token, err = vcapool.NewAuthToken(&user); err != nil {
		return vcago.NewStatusInternal(err)
	}
	c.SetCookie(token.AccessCookie())
	c.SetCookie(token.RefreshCookie())
	return c.JSON(vcago.NewResponse("refresh_token", user).Selected())
}

func LogoutHandler(c echo.Context) (err error) {
	user := new(vcapool.AccessToken)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	c.SetCookie(vcapool.ResetAccessCookie())
	c.SetCookie(vcapool.ResetRefreshCookie())
	return c.JSON(vcago.NewResponse("logout", user.ID).Executed())
}
