package token

import (
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type AssetsHandler struct {
	vcago.Handler
}

var Assets = &AssetsHandler{*vcago.NewHandler("assets")}

func (i *AssetsHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.GET("/:id", i.GetByID, accessCookie)
}

func (i *AssetsHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AssetsID)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	return c.File("/public/files/" + body.ID)
}
