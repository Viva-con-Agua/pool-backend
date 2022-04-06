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
	if !token.Roles.Validate("employee;admin") {
		body.UserID = token.ID
	}
	result := new(vcapool.Address)
	if result, err = body.Get(c.Ctx()); err != nil {
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
	if token.AddressID != body.ID {
		return vcago.NewPermissionDenied("address", body.ID)
	}
	result := new(vcapool.Address)
	if result, err = body.Update(c.Ctx()); err != nil {
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
	if token.AddressID != body.ID {
		return vcago.NewPermissionDenied("address", body.ID)
	}
	if err = body.Delete(c.Ctx()); err != nil {
		return
	}
	return c.Deleted(body.ID)
}

//TODO
func AddressList(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(dao.AddressQuery)
	if err = vcago.BindAndValidate(c, body); err != nil {
		return
	}
	result := new(dao.AddressList)
	if err = result.Get(ctx, body.Filter()); err != nil {
		return
	}
	return vcago.NewSelected("address_list", result)
}
