package token

import (
	"log"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type AddressHandler struct {
	vcago.Handler
}

var Address = &AddressHandler{*vcago.NewHandler("address")}

func (i *AddressHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)
}

func (i *AddressHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AddressCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := body.Address(token.ID)
	if err = dao.AddressesCollection.InsertOne(c.Ctx(), result); err != nil {
		return
	}
	if err = dao.IDjango.Post(result, "/v1/pool/address/"); err != nil {
		log.Print(err)
	}
	return c.Created(result)
}

func (i *AddressHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AddressParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Address)
	if err = dao.AddressesCollection.FindOne(c.Ctx(), body.Filter(token), result); err != nil {
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
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Address)
	if err = dao.AddressesCollection.UpdateOne(c.Ctx(), body.Filter(token), vmdb.UpdateSet(body), result); err != nil {
		return
	}
	if err = dao.IDjango.Post(result, "/v1/pool/address/"); err != nil {
		log.Print(err)
	}
	return c.Updated(result)
}

func (i *AddressHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AddressParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.AddressesCollection.DeleteOne(c.Ctx(), body.Filter(token)); err != nil {
		return
	}
	var result *models.NVM
	if result, err = dao.NVMWithdraw(c.Ctx(), token); err != nil {
		return
	}
	if err = dao.IDjango.Post(&models.Address{UserID: token.ID}, "/v1/pool/address/"); err != nil {
		log.Print(err)
	}
	if err = dao.IDjango.Post(result, "/v1/pool/profile/nvm/"); err != nil {
		log.Print(err)
	}
	return c.Deleted(body.ID)
}

// TODO
func (i *AddressHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AddressQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new([]models.Address)
	if err = dao.AddressesCollection.Find(c.Ctx(), body.Filter(token), result); err != nil {
		return
	}
	return c.Listed(result)
}
