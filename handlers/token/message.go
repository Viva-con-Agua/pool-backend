package token

import (
	"log"
	"net/http"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
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
	crew := new(models.Crew)
	if err = dao.CrewsCollection.FindOne(c.Ctx(), bson.D{{Key: "_id", Value: token.CrewID}}, crew); err != nil {
		log.Print("No crew for user")
	}

	result := new(models.Message)
	if err = dao.MessageCollection.FindOne(c.Ctx(), body.Filter(token, crew), result); err != nil {
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
	crew := new(models.Crew)
	if err = dao.CrewsCollection.FindOne(c.Ctx(), bson.D{{Key: "_id", Value: token.CrewID}}, crew); err != nil {
		log.Print("No crew for user")
	}
	result := new(models.Message)
	if err = dao.MessageCollection.UpdateOne(
		c.Ctx(),
		body.Filter(token, crew),
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
	crew := new(models.Crew)
	if err = dao.CrewsCollection.FindOne(c.Ctx(), bson.D{{Key: "_id", Value: token.CrewID}}, crew); err != nil {
		log.Print("No crew for user")
	}
	if err = dao.MessageCollection.DeleteOne(c.Ctx(), body.Filter(token, crew)); err != nil {
		return
	}
	return c.Deleted(body.ID)
}

func (i *MessageHandler) SendCycular(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(vmod.IDParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	var result *models.Message
	var mail *vcago.CycularMail
	if result, mail, err = dao.MessageSendCycular(c.Ctx(), body, token); err != nil {
		return
	}
	vcago.Nats.Publish("system.mail.cycular", mail)
	return c.SuccessResponse(http.StatusOK, "send mail", "message", result)
}
