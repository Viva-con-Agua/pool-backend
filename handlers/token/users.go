package token

import (
	"log"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	vcago.Handler
}

var User = &UserHandler{*vcago.NewHandler("user")}

func (i *UserHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.GET("", i.Get, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)
}

func (i *UserHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.UserDelete(c.Ctx(), body.ID); err != nil {
		return
	}
	if err = dao.IDjango.Post(&models.Profile{UserID: body.ID}, "/v1/pool/profile/delete/"); err != nil {
		log.Print(err)
	}
	return c.Deleted(body.ID)
}

func (i *UserHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.User)
	if err = dao.UserCollection.AggregateOne(c.Ctx(), body.Pipeline(), result); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *UserHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	var result *[]models.User
	if result, err = dao.UsersGet(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}
