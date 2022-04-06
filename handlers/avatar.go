package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type AvatarHandler struct {
	vcago.Handler
}

func NewAvatarHandler() *AvatarHandler {
	handler := vcago.NewHandler("avatar")
	return &AvatarHandler{
		*handler,
	}
}

func (i *AvatarHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(dao.AvatarCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	body.UserID = token.ID
	result := new(vcapool.Avatar)
	if result, err = body.Create(c.Ctx()); err != nil {
		return
	}
	return c.Created(result)
}

func (i *AvatarHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(dao.AvatarUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if token.AvatarID != body.ID {
		return vcago.NewPermissionDenied("avatar", body.ID)
	}
	result := new(vcapool.Avatar)
	if result, err = body.Update(c.Ctx()); err != nil {
		return
	}
	return c.Updated(result)

}

func (i *AvatarHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(dao.Avatar)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	id := c.Param("c")
	if token.AvatarID != id {
		return vcago.NewPermissionDenied("avatar", body.ID)
	}
	if err = body.Delete(c.Ctx(), id); err != nil {
		return
	}
	return c.Deleted(id)
}
