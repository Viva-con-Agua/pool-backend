package token

import (
	"net/http"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type MessageHandler struct {
	vcago.Handler
}

var Message = &MessageHandler{*vcago.NewHandler("message")}

func (i *MessageHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.GET("/send_cycular/:id", i.SendCycular, accessCookie)
	group.POST("", i.Create, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)
}

// Create
// @Security CookieAuth
// @Summary Create a Message
// @Description creates an Message object.
// @Tags /mails/message
// @Accept json
// @Produce json
// @Param form body models.MessageCreate true "Message Data"
// @Model: vcago.Response
// @Success 201 {object} vcago.ResponseCreated{payload=models.Message} "Message successfully created"
// @Failure 400 {object} vcago.BindErrorResponse{} "Bind Error"
// @Failure 400 {object} vcago.ValidationErrorResponse{} "Validation Error"
// @Failure 409 {object} vcago.MongoDuplicatedErrorResponse{} "Duplicated Key"
// @Router /mails/message [post]
func (i *MessageHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.MessageCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Message)
	if result, err = dao.MessageInsert(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Created(result)
}

// GetByID
// @Security CookieAuth
// @Summary Get a Message by ID
// @Tags /mails/message
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Model: vcago.Response
// @Success 200 {object} vcago.ResponseSelected{payload=models.Message} "Message successfully selected"
// @Failure 404 {object} vcago.MongoNoDocumentErrorResponse{} "No Document with given ID"
// @Router /mails/message/{id} [get]
func (i *MessageHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.MessageParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Message)
	if result, err = dao.MessageGetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

// Update
// @Security CookieAuth
// @Summary Update a Message
// @Tags /mails/message
// @Accept json
// @Produce json
// @Param form body models.DepositUpdate true "Message Data"
// @Model: vcago.Response
// @Success 200 {object} vcago.ResponseUpdated{payload=models.Message} "Message successfully updated"
// @Failure 400 {object} vcago.BindErrorResponse{} "Bind Error"
// @Failure 400 {object} vcago.ValidationErrorResponse{} "Validation Error"
// @Failure 404 {object} vcago.MongoNoDocumentErrorResponse{} "No Document with given ID"
// @Router /mails/message [put]
func (i *MessageHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.MessageUpdate)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Message)
	if result, err = dao.MessageUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Updated(result)
}

// DeleteByID
// @Security CookieAuth
// @Summary Delete a Message by ID
// @Tags /mails/message
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} vcago.ResponseDeleted{payload=string} "Message successfully deleted"
// @Failure 404 {object} vcago.MongoNoDocumentErrorResponse{} "No Document with given ID"
// @Router /mails/message/{id} [delete]
func (i *MessageHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.MessageParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.MessageDelete(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Deleted(body.ID)
}

// SendCycular
// @Security CookieAuth
// @Summary Get a Message by ID
// @Tags /mails/message
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Model: vcago.Response
// @Success 200 {object} vcago.ResponseSelected{payload=models.Message} "Message successfully sended"
// @Failure 404 {object} vcago.MongoNoDocumentErrorResponse{} "No Document with given ID"
// @Router /mails/message/send_cycular/{id} [get]
func (i *MessageHandler) SendCycular(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.MessageParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Message)
	mail := new(vcago.CycularMail)
	if result, mail, err = dao.MessageSendCycular(c.Ctx(), body, token); err != nil {
		return
	}
	vcago.Nats.Publish("system.mail.cycular", mail)
	return c.SuccessResponse(http.StatusOK, "send mail", "message", result)
}
