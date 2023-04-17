package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type MailboxHandler struct {
	vcago.Handler
}

var Mailbox = &MailboxHandler{*vcago.NewHandler("mailbox")}

func (i *MailboxHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.GET("/:id", i.GetByID, accessCookie)
}

func (i *MailboxHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.MailboxParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Mailbox)
	if result, err = dao.MailboxGetByID(c.Ctx(), body.ID, token); err != nil {
		return
	}
	return c.Selected(result)
}
