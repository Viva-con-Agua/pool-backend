package handlers

import (
	"log"
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
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
	err = userDAO.Get(ctx, tokenUser.ID)
	log.Print(err)
	if vcago.MongoNoDocuments(err) {
		if userDAO, err = dao.CreateUserFromToken(ctx, tokenUser); err != nil {
			return
		}
	}
	if err != nil {
		return
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
	if err = userDAO.Get(ctx, userID); err != nil {
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
	return
}
