package token

import (
	"log"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type AddressHandler struct {
	vcago.Handler
}

var Address = &AddressHandler{*vcago.NewHandler("address")}

func (i *AddressHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.POST("/create", i.UsersCreate, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.PUT("/update", i.UsersUpdate, accessCookie)
	group.DELETE("/delete/:id", i.UsersDelete, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)
}

func (i *AddressHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AddressCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Address)
	if result, err = dao.AddressInsert(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/address/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Created(result)
}

func (i *AddressHandler) UsersCreate(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UsersAddressCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Address)
	if result, err = dao.UsersAddressInsert(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/address/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Created(result)
}

func (i *AddressHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AddressQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new([]models.Address)
	if result, err = dao.AddressGet(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *AddressHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AddressParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Address)
	if result, err = dao.AddressGetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *AddressHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AddressUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Address)
	if result, err = dao.AddressUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/address/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Updated(result)
}

func (i *AddressHandler) UsersUpdate(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AddressUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Address)
	if result, err = dao.UsersAddressUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/address/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Updated(result)
}

func (i *AddressHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AddressParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.NVM)
	if result, err = dao.AddressDelete(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(&models.Address{UserID: token.ID}, "/v1/pool/address/"); err != nil {
			log.Print(err)
		}
		if err = dao.IDjango.Post(result, "/v1/pool/profile/nvm/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Deleted(body.ID)
}

func (i *AddressHandler) UsersDelete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AddressParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.NVM)
	if result, err = dao.UsersAddressDelete(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(&models.Address{UserID: token.ID}, "/v1/pool/address/"); err != nil {
			log.Print(err)
		}
		if err = dao.IDjango.Post(result, "/v1/pool/profile/nvm/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Deleted(body.ID)
}
