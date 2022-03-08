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
	user := new(vcapool.User)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	if user.ID != body.UserID {
		return vcago.NewStatusPermissionDenied()
	}
	if err = body.Update(ctx); err != nil {
		return
	}
	//active state
	result := new(dao.UserActive)
	if err = result.Get(ctx, bson.M{"user_id": user.ID}); err != nil {
		return
	}
	if result, err = result.Withdraw(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("user_crew", body).Updated())
}

func DeleteUserCrew(c echo.Context) (err error) {
	ctx := c.Request().Context()
	user := new(vcapool.User)
	if user, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	result := new(dao.UserCrew)
	if err = result.Delete(ctx, bson.M{"user_id": user.ID}); err != nil {
		return
	}
	//active state
	resultA := new(dao.UserActive)
	if err = resultA.Get(ctx, bson.M{"user_id": user.ID}); err != nil {
		return
	}
	if resultA, err = resultA.Withdraw(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("user_crew", result).Deleted())

}
