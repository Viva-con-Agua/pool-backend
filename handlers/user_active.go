package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func UserActiveRequest(c echo.Context) (err error) {
	ctx := c.Request().Context()
	user := new(vcapool.AccessToken)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	result := new(dao.UserActive)
	if err = result.Get(ctx, bson.M{"user_id": user.ID}); err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	if user.CrewID == "" {
		return vcago.NewBadRequest("user_active", "not an crew member")
	}
	if vcago.MongoNoDocuments(err) {
		err = nil
		if result, err = result.Create(ctx, user.ID); err != nil {
			return
		}
	} else {
		if err = result.Request(ctx); err != nil {
			return
		}
	}
	return vcago.NewCreated("user_active", result)
}

//ConfirmUserActive is the webapp handler for confirm the active state of an user.
func UserActiveConfirm(c echo.Context) (err error) {
	ctx := c.Request().Context()
	//validate and bind body
	body := new(dao.UserActiveRequest)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	//get requested user from token
	userReq := new(vcapool.AccessToken)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	//check if requested user has the network or operation permission
	if !userReq.PoolRoles.Validate("network;operation") && !userReq.Roles.Validate("employee") {
		return vcago.NewPermissionDenied("user_active")
	}
	//check if requested user and the target users has the same crew
	if !userReq.Roles.Validate("employee") {
		userCrew := new(dao.UserCrew)
		if err = userCrew.Permission(ctx, bson.M{"user_id": body.UserID, "crew_id": userReq.CrewID}); err != nil {
			return
		}
	}
	//confirm active state
	result := new(dao.UserActive)
	if err = result.Confirm(ctx, body.UserID); err != nil {
		return
	}
	mailData := new(vcago.MailData)
	if mailData, err = dao.GetSendMail(ctx, userReq.ID, result.UserID, "active_confirmed"); err != nil {
		return
	}
	vcago.Nats.Publish("mail.send", mailData)
	dao.MailSend.Send(mailData)
	//response the result as vcago.Response
	return vcago.NewExecuted("user_active", result)
}

func UserActiveReject(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.UserActiveRequest)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	//get requested user from token
	userReq := new(vcapool.AccessToken)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	//check if requested user has the network or operation permission
	if !userReq.PoolRoles.Validate("network;operation") && !userReq.Roles.Validate("employee") {
		return vcago.NewPermissionDenied("user_active")
	}
	//check if requested user and the target users has the same crew
	if !userReq.Roles.Validate("employee") {
		userCrew := new(dao.UserCrew)
		if err = userCrew.Permission(ctx, bson.M{"user_id": body.UserID, "crew_id": userReq.CrewID}); err != nil {
			return
		}
	}
	result := new(dao.UserActive)
	if err = result.Reject(ctx, body.UserID); err != nil {
		return
	}
	result2 := new(dao.UserNVM)
	err = result2.Get(ctx, bson.M{"user_id": body.UserID})
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	if !vcago.MongoNoDocuments(err) {
		err = nil
		if result2, err = result2.Reject(ctx, body.UserID); err != nil {
			return
		}
	}
	mailData := new(vcago.MailData)
	if mailData, err = dao.GetSendMail(ctx, userReq.ID, result.UserID, "active_rejected"); err != nil {
		return
	}
	dao.MailSend.Send(mailData)
	return vcago.NewExecuted("user_active", result)
}

func UserActiveWithdraw(c echo.Context) (err error) {
	ctx := c.Request().Context()
	user := new(vcapool.AccessToken)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	result := new(dao.UserActive)
	if err = result.Get(ctx, bson.M{"user_id": user.ID}); err != nil {
		return
	}
	if err = result.Withdraw(ctx); err != nil {
		return
	}
	result2 := new(dao.UserNVM)
	err = result2.Get(ctx, bson.M{"user_id": user.ID})
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	if !vcago.MongoNoDocuments(err) {
		err = nil
		if result2, err = result2.Withdraw(ctx); err != nil {
			return
		}
	}
	return vcago.NewExecuted("user_active", result)
}
