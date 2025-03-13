package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type OrganisationHandler struct {
	vcago.Handler
}

var Organisation = &OrganisationHandler{*vcago.NewHandler("organisation")}

func (i *OrganisationHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)
}

func (i *OrganisationHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.OrganisationCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Organisation)
	if result, err = dao.OrganisationInsert(c.Ctx(), body, token); err != nil {
		return
	}
	//go func() {
	//	if err = dao.IDjango.Post(result, "/v1/pool/organisation/"); err != nil {
	//		log.Print(err)
	//	}
	//}()
	return c.Created(result)
}

func (i *OrganisationHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.OrganisationQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new([]models.Organisation)
	if result, err = dao.OrganisationGet(c.Ctx(), body); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *OrganisationHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.OrganisationParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new(models.Organisation)
	if result, err = dao.OrganisationGetByID(c.Ctx(), body); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *OrganisationHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.OrganisationUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Organisation)
	if result, err = dao.OrganisationUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	//go func() {
	//	if err = dao.IDjango.Post(body, "/v1/pool/organisation/"); err != nil {
	//		log.Print(err)
	//	}
	//}()
	return c.Updated(result)
}

func (i *OrganisationHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.OrganisationParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.OrganisationDelete(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
