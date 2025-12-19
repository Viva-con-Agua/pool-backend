package token

import (
	"net/http"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type ProfileHandler struct {
	vcago.Handler
}

var Profile = &ProfileHandler{*vcago.NewHandler("profile")}

func (i *ProfileHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.GET("/sync/:id", i.UserSync, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.PUT("/update", i.UsersUpdate, accessCookie)
}

func (i *ProfileHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ProfileCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Profile)
	if result, err = dao.ProfileInsert(c.Ctx(), body, token); err != nil {
		return
	}
	dao.ProfileSync(result)
	return c.Created(result)
}

func (i *ProfileHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ProfileUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Profile)
	if result, err = dao.ProfileUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	dao.ProfileSync(result)
	return c.Updated(result)
}

func (i *ProfileHandler) UserSync(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ProfileParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = body.ProfileSyncPermission(token); err != nil {
		return
	}
	result := new(models.User)
	if result, err = dao.UserSync(c.Ctx(), body, token); err != nil {
		return
	}
	if _, err = dao.ProfileSync(&result.Profile); err != nil {
		return
	}
	if _, err = dao.NewsletterSync(result); err != nil {
		return
	}
	return c.SuccessResponse(http.StatusOK, "successfully_synced", "active", nil)
}

func (i *ProfileHandler) UsersUpdate(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ProfileUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Profile)
	if result, err = dao.UsersProfileUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	dao.ProfileSync(result)
	return c.Updated(result)
}
