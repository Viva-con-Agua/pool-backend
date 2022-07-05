package token

import (
	"log"
	"net/http"
	"pool-user/dao"
	"pool-user/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//ActiveHandler represents the vcago.Handler for the Active model.
type ActiveHandler struct {
	vcago.Handler
}

//Active is an ActiveHandler.
var Active = &ActiveHandler{*vcago.NewHandler("active")}

func (i *ActiveHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.GET("/request", i.Request, vcapool.AccessCookieConfig())
	group.POST("/confirm", i.Confirm, vcapool.AccessCookieConfig())
	group.POST("/reject", i.Reject, vcapool.AccessCookieConfig())
	group.GET("/withdraw", i.Withdraw, vcapool.AccessCookieConfig())
}

//Request handles an active request call.
func (i *ActiveHandler) Request(cc echo.Context) (err error) {
	//load context
	c := cc.(vcago.Context)
	//get access token
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	//check permissions for active request
	if err = models.ActiveRequestPermission(token); err != nil {
		return
	}
	//update active model into database. For filter the user_id key it the value ot the access token ID.
	result := new(models.Active)
	if err = dao.ActiveCollection.UpdateOne(
		c.Ctx(),
		bson.D{{Key: "user_id", Value: token.ID}},
		vmdb.NewUpdateSet(models.ActiveRequest()),
		result,
	); err != nil {
		log.Print(err)
		return
	}
	//return "successfully requested."
	return c.SuccessResponse(http.StatusOK, "successfully_requested", "active", result)
}

//Confirm handles an active confirm request.
func (i *ActiveHandler) Confirm(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	//validate and bind body
	body := new(models.ActiveParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	//get requested user from token
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	//check permissions for update an other users active model.
	if err = models.ActivePermission(token); err != nil {
		return
	}
	//update active model.
	result := new(models.Active)
	if err = dao.ActiveCollection.UpdateOne(
		c.Ctx(),
		body.Filter(token),
		vmdb.NewUpdateSet(models.ActiveConfirm()),
		result,
	); err != nil {
		return
	}
	//confirm active state
	/*result := new(vcapool.UserActive)
	if result, err = body.Confirm(ctx, body.UserID); err != nil {
		return
	}
	mailData := new(vcago.MailData)
	if mailData, err = dao.GetSendMail(ctx, userReq.ID, result.UserID, "active_confirmed"); err != nil {
		return
	}
	vcago.Nats.Publish("mail.send", mailData)
	dao.MailSend.Send(mailData)
	//response the result as vcago.Response
	*/
	return c.SuccessResponse(http.StatusOK, "successfully_confirmed", "active", result)
}

func (i *ActiveHandler) Reject(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ActiveParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	//get requested user from token
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	//check permissions for update an other users active model.
	if err = models.ActivePermission(token); err != nil {
		return
	}
	//update active model.
	result := new(models.Active)
	if err = dao.ActiveCollection.UpdateOne(
		c.Ctx(),
		body.Filter(token),
		vmdb.NewUpdateSet(models.ActiveReject()),
		result,
	); err != nil {
		return
	}
	//reject nvm state
	if err = dao.NVMCollection.UpdateOne(
		c.Ctx(),
		bson.D{{Key: "user_id", Value: body.UserID}},
		vmdb.NewUpdateSet(models.NVMReject()),
		nil,
	); err != nil && err != mongo.ErrNoDocuments {
		return
	}
	/*
		mailData := new(vcago.MailData)
		if mailData, err = dao.GetSendMail(ctx, userReq.ID, result.UserID, "active_rejected"); err != nil {
			return
		}
		dao.MailSend.Send(mailData)*/
	return c.SuccessResponse(http.StatusOK, "successfully_rejected", "active", result)
}

func (i *ActiveHandler) Withdraw(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	//reject nvm state
	result := new(models.Active)
	if err = dao.ActiveCollection.UpdateOne(
		c.Ctx(),
		bson.D{{Key: "user_id", Value: token.ID}},
		vmdb.NewUpdateSet(models.NVMWithdraw()),
		result,
	); err != nil {
		return
	}
	//reject nvm state
	if err = dao.NVMCollection.UpdateOne(
		c.Ctx(),
		bson.D{{Key: "user_id", Value: token.ID}},
		vmdb.NewUpdateSet(models.NVMReject()),
		nil,
	); err != nil && vmdb.ErrNoDocuments(err) {
		return
	}
	/*
		mailData := new(vcago.MailData)
		if mailData, err = dao.GetSendMail(ctx, userReq.ID, result.UserID, "active_rejected"); err != nil {
			return
		}
		dao.MailSend.Send(mailData)*/
	return c.SuccessResponse(http.StatusOK, "successfully_rejected", "active", result)
}
