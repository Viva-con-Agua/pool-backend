package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateUserCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.UserCrewCreateRequest)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	user := new(vcapool.User)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	result := new(dao.UserCrew)
	if result, err = body.Create(ctx, user.ID); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("user_crew", result).Created())
}

func UpdateUserCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.UserCrew)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userReq := new(vcapool.User)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	if userReq.ID != body.UserID {
		return vcago.NewStatusPermissionDenied()
	}
	if err = body.Update(ctx); err != nil {
		return
	}
	//active state
	result := new(dao.UserActive)
	err = result.Get(ctx, bson.M{"user_id": userReq.ID})
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	if !vcago.MongoNoDocuments(err) {
		err = nil
		if result, err = result.Withdraw(ctx); err != nil {
			return
		}
	}
	result2 := new(dao.UserNVM)
	err = result2.Get(ctx, bson.M{"user_id": userReq.ID})
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	if !vcago.MongoNoDocuments(err) {
		err = nil
		if result2, err = result2.Withdraw(ctx); err != nil {
			return
		}
	}
	return c.JSON(vcago.NewResponse("user_crew", body).Updated())
}

func DeleteUserCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	userReq := new(vcapool.User)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	result := new(dao.UserCrew)
	if err = result.Delete(ctx, bson.M{"user_id": userReq.ID}); err != nil {
		return
	}
	//active state
	//active state
	resultA := new(dao.UserActive)
	err = resultA.Get(ctx, bson.M{"user_id": userReq.ID})
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	if !vcago.MongoNoDocuments(err) {
		err = nil
		if resultA, err = resultA.Withdraw(ctx); err != nil {
			return
		}
	}
	result2 := new(dao.UserNVM)
	err = result2.Get(ctx, bson.M{"user_id": userReq.ID})
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	if !vcago.MongoNoDocuments(err) {
		err = nil
		if result2, err = result2.Withdraw(ctx); err != nil {
			return
		}
	}
	return c.JSON(vcago.NewResponse("user_crew", result).Deleted())

}
