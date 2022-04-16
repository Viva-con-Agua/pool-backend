package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type AddressHandler struct {
	vcago.Handler
}

func NewAddressHandler() *AddressHandler {
	handler := vcago.NewHandler("address")
	return &AddressHandler{
		*handler,
	}
}

func (i *AddressHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(dao.AddressCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(vcapool.Address)
	if result, err = body.Create(c.Ctx(), token); err != nil {
		return
	}
	return c.Created(result)
}

func (i *AddressHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(dao.AddressParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(vcapool.Address)
	if result, err = body.Get(c.Ctx(), token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *AddressHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(dao.AddressUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(vcapool.Address)
	if result, err = body.Update(c.Ctx(), token); err != nil {
		return
	}
	return c.Updated(result)
}

func (i *AddressHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(dao.AddressParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = body.Delete(c.Ctx(), token); err != nil {
		return
	}
	return c.Deleted(body.ID)
}

//TODO
func (i *AddressHandler) List(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(dao.AddressQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(vcapool.AddressList)
	if result, err = body.List(c.Ctx(), token); err != nil {
		return
	}
	return c.Listed(result)
}
