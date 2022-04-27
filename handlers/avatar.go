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
	result := new(vcapool.Avatar)
	if result, err = body.Create(c.Ctx(), token); err != nil {
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
	result := new(vcapool.Avatar)
	if result, err = body.Update(c.Ctx(), token); err != nil {
		return
	}
	return c.Updated(result)

}

func (i *AvatarHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(dao.AvatarParam)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = body.Delete(c.Ctx(), token); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
