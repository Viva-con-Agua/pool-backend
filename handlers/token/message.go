package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type MessageHandler struct {
	vcago.Handler
}

var Message = &MessageHandler{*vcago.NewHandler("message")}

func (i *MessageHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	//group.GET("/send_cycular/:id", i.SendCycular, accessCookie)
	group.POST("", i.Create, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)
}

func (i *MessageHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.MessageCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := body.MessageSub(token)
	if err = dao.MessageCollection.InsertOne(c.Ctx(), result); err != nil {
		return
	}
	return c.Created(result)
}

func (i *MessageHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.MessageParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}

	result := new(models.Message)
	if err = dao.MessageCollection.FindOne(c.Ctx(), body.Filter(token), result); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *MessageHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.MessageUpdate)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Message)
	if err = dao.MessageCollection.UpdateOne(
		c.Ctx(),
		body.Filter(token),
		vmdb.UpdateSet(body),
		result,
	); err != nil {
		return
	}
	return c.Updated(body)
}

func (i *MessageHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.MessageParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.MessageCollection.DeleteOne(c.Ctx(), body.Filter(token)); err != nil {
		return
	}
	return c.Deleted(body.ID)
}

/*
func (i *MessageHandler) SendCycular(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.MessageParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Message)
	if err = dao.MessageCollection.FindOne(c.Ctx(), body.Filter(token), result); err != nil {
		return
	}
	if result.RecipientGroup.Type == "crew" {
		userList := new([]models.User)
		if userList, err = dao.UserGetRequest(&result.RecipientGroup); err != nil {
			return
		}
		result.To = *userList
	} else if result.RecipientGroup.Type == "event" {
		return vcago.NewBadRequest("message", "event is currently not supported", result.RecipientGroup)
	} else {
		return vcago.NewBadRequest("message", "type is not supported", result.RecipientGroup)
	}

	mail := vcago.NewCycularMail(result.From, result.To.Emails(), result.Subject, result.Message)
	dao.MailSend.PostCycularMail(mail)
	if err = dao.MessageCollection.InsertMany(c.Ctx(), *result.Inbox()); err != nil {
		return
	}
	result.Type = "outbox"
	if err = dao.MessageCollection.UpdateOne(
		c.Ctx(),
		bson.D{{Key: "_id", Value: result.ID}},
		vmdb.UpdateSet(result.MessageUpdate()),
		result,
	); err != nil {
		return
	}
	return c.SuccessResponse(http.StatusOK, "send mail", "message", result)
}*/
