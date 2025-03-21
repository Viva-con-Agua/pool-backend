package token

import (
	"log"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type CrewHandler struct {
	vcago.Handler
}

var Crew = &CrewHandler{*vcago.NewHandler("crew")}

func (i *CrewHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.GET("/own", i.GetAsMember, accessCookie)
	group.GET("/public", i.GetPublic)
	group.GET("/:id", i.GetByID, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)
}

// Create
// @Security CookieAuth
// @Summary Create a Crew
// @Description creates an  Crew object.
// @Tags /crews
// @Accept json
// @Produce json
// @Param form body models.CrewCreate true "Crew Data"
// @Model: vcago.Response
// @Success 201 {object} vcago.Response{payload=models.Crew}
// @Router /crews [post]
func (i *CrewHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Crew)
	if result, err = dao.CrewInsert(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/crew/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Created(result)
}

// Get
// @Security CookieAuth
// @Summary Get a List of  Crews
// @Tags /crews
// @Accept json
// @Produce json
// @Param   q query   models.CrewQuery   false  "string collection"  collectionFormat(multi)
// @Model: vcago.Response
// @Success 200 {object} vcago.Response{payload=[]models.Crew}
// @Failure 400 {object} vcago.Response{}
// @Router /crews [get]
func (i *CrewHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new([]models.Crew)
	if result, err = dao.CrewGet(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

// GetPublic
// @Summary Get a List of CrewPublic
// @Tags /crews
// @Accept json
// @Produce json
// @Param   q query   models.CrewQuery   false  "string collection"  collectionFormat(multi)
// @Model: vcago.Response
// @Success 200 {object} vcago.Response{payload=[]models.CrewPublic}
// @Failure 400 {object} vcago.Response{}
// @Router /crews/public [get]
func (i *CrewHandler) GetPublic(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new([]models.CrewPublic)
	if result, err = dao.CrewPublicGet(c.Ctx(), body); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *CrewHandler) GetAsMember(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Crew)
	if result, err = dao.CrewGetAsMember(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

// GetByID
// @Security CookieAuth
// @Summary Get a Crew by ID
// @Tags /crews
// @Accept json
// @Produce json
// @Param id path string true "Crew ID"
// @Model: vcago.Response
// @Success 200 {object} vcago.Response{payload=models.Crew}
// @Router /users/address/{id} [get]
func (i *CrewHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Crew)
	if result, err = dao.CrewGetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

// Update
// @Security CookieAuth
// @Summary Update a Crew
// @Tags /crews
// @Accept json
// @Produce json
// @Param form body models.CrewUpdate true "Crew Data"
// @Model: vcago.Response
// @Success 200 {object} vcago.Response{payload=models.Crew}
// @Router /crews [put]
func (i *CrewHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Crew)
	if result, err = dao.CrewUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(body, "/v1/pool/crew/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Updated(result)
}

// DeleteByID
// @Security CookieAuth
// @Summary Delete a Crew by ID
// @Tags /crews
// @Accept json
// @Produce json
// @Param id path string true "Crew ID"
// @Success 200 {object} vmod.DeletedResponse
// @Router /crews/{id} [delete]
func (i *CrewHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.CrewDelete(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
