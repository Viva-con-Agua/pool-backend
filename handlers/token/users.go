package token

import (
	"log"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
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
	group.PUT("/organisation", i.UpdateOrganisation, accessCookie)
	group.GET("/crew", i.GetUsersByCrew, accessCookie)
	group.GET("/crew/public", i.GetMinimal, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)
	group.GET("/api_key/:id", i.GetByIDApiKey, vcago.KeyAuthMiddleware())
}

func (i *UserHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new([]models.ListUser)
	var listSize int64
	if result, listSize, err = dao.UsersGet(body, token); err != nil {
		return
	}
	return c.Listed(result, listSize)
}

func (i *UserHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.User)
	if result, err = dao.UsersUserGetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *UserHandler) GetByIDApiKey(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new(models.User)
	if result, err = dao.UsersUserGetByIDAPIKey(c.Ctx(), body); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *UserHandler) GetUsersByCrew(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new([]models.UserBasic)
	if result, err = dao.UsersGetByCrew(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *UserHandler) GetMinimal(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new([]models.UserMinimal)
	if result, err = dao.UsersMinimalGet(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *UserHandler) UpdateOrganisation(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserOrganisationUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.User)
	if result, err = dao.UserOrganisationUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Updated(result)
}

func (i *UserHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.UsersDeleteUser(c.Ctx(), body, token); err != nil {
		return
	}
	vcago.Nats.Publish("user.delete", body)
	go func() {
		if err = dao.IDjango.Post(&models.Profile{UserID: body.ID}, "/v1/pool/profile/delete/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Deleted(body.ID)
}
