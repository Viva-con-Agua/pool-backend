package token

import (
	"log"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmod"
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

// Create
// @Security CookieAuth
// @Summary Create a Address
// @Description creates an  Address object.
// @Tags Address
// @Accept json
// @Produce json
// @Param form body models.AddressCreate true "Address Data"
// @Model: vcago.Response
// @Success 201 {object} vcago.Response{payload=models.Address}
// @Router /users/address [post]
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

// UsersCreate
// TODO: delete
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

// Get
// @Security CookieAuth
// @Summary Get a List of  Addresss
// @Tags Address
// @Accept json
// @Produce json
// @Param   q query   models.AddressQuery   false  "string collection"  collectionFormat(multi)
// @Model: vcago.Response
// @Success 200 {object} vcago.Response{payload=[]models.Address}
// @Failure 400 {object} vcago.Response{}
// @Router /users/address [get]
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

// GetByID
// @Security CookieAuth
// @Summary Get a  Address by ID
// @Tags Address
// @Accept json
// @Produce json
// @Param id path string true "Address ID"
// @Model: vcago.Response
// @Success 200 {object} vcago.Response{payload=models.Address}
// @Router /users/address/{id} [get]
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

// Update
// @Security CookieAuth
// @Summary Get a  Address by ID
// @Tags Address
// @Accept json
// @Produce json
// @Param form body models.AddressUpdate true "Address Data"
// @Model: vcago.Response
// @Success 200 {object} vcago.Response{payload=models.Address}
// @Router /users/address [put]
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

// UsersUpdate
// TODO: delete
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

// DeleteByID
// @Security CookieAuth
// @Summary Delete a  Address by ID
// @Tags Address
// @Accept json
// @Produce json
// @Param id path string true "Address ID"
// @Success 200 {object} vmod.DeletedResponse
// @Router /users/address/{id} [delete]
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
	var result *vmod.DeletedResponse
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

// UsersDelete
// TODO: delete
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
	var result *vmod.DeletedResponse
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
