package token

import (
	"log"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
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

func (i *CrewHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = models.CrewPermission(token); err != nil {
		return
	}
	result := body.Crew()
	if result, err = dao.CrewInsert(c.Ctx(), body); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/crew/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Created(result)
}

func (i *CrewHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = models.CrewPermission(token); err != nil {
		return
	}
	result := new([]models.Crew)
	if dao.CrewsCollection.Find(c.Ctx(), body.Filter(), result); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *CrewHandler) GetPublic(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	body.Status = "active"
	result := new([]models.CrewPublic)
	if dao.CrewsCollection.Find(c.Ctx(), body.Filter(), result); err != nil {
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
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	body.Status = "active"
	body.ID = []string{token.CrewID}
	result := new(models.Crew)
	if dao.CrewsCollection.FindOne(c.Ctx(), body.Filter(), result); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *CrewHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Crew)
	if result, err = dao.CrewGetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *CrewHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = models.CrewUpdatePermission(token); err != nil {
		return
	}
	if _, err = dao.CrewUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(body, "/v1/pool/crew/"); err != nil {
			log.Print(err)
		}
	}()
	return vcago.NewUpdated("crew", body)
}

func (i *CrewHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.CrewParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = models.CrewPermission(token); err != nil {
		return
	}
	if err = dao.CrewDelete(c.Ctx(), body); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
