package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
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

// GetByID
// @Security CookieAuth
// @Summary Get a Mailbox by ID
// @Tags /mails/mailbox
// @Accept json
// @Produce json
// @Param id path string true "Mailbox ID"
// @Model: vcago.Response
// @Success 200 {object} vcago.ResponseSelected{payload=models.Mailbox}
// @Failure 404 {object} vcago.MongoNoDocumentErrorResponse{} "No Document with given ID"
// @Router /mails/mailbox/{id} [get]
func (i *MailboxHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.MailboxParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Mailbox)
	if result, err = dao.MailboxGetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}
