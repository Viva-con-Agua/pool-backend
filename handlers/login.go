package handlers

import (
	"net/http"
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
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
	user := new(dao.User)
	err = user.Get(ctx, tokenUser.ID)
	if vcago.MongoNoDocuments(err) {
		user = dao.NewUser(tokenUser)
		err = user.Create(ctx)
	}
	if err != nil {
		return
	}
	token := new(dao.AuthToken)
	if token, err = user.ToAuthToken(); err != nil {
		return vcago.NewStatusInternal(err)
	}
	c.SetCookie(token.Token.AccessCookie())
	c.SetCookie(token.Token.RefreshCookie())
	return c.JSON(http.StatusOK, user.ToVPUser())
}
