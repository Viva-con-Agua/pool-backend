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

type EventHandler struct {
	vcago.Handler
}

var Event = &EventHandler{*vcago.NewHandler("event")}

func (i *EventHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.GET("/public", i.GetPublic)
	group.GET("", i.Get, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)

}

func (i *EventHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	database := body.EventDatabase(token)
	result := new(models.Event)
	if result, err = dao.EventInsert(c.Ctx(), database, token); err != nil {
		return
	}
	result.EditorID = token.ID
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/event/create/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Created(result)
}

func (i *EventHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Event)
	if err = dao.EventCollection.AggregateOne(
		c.Ctx(),
		models.EventPipeline(token).Match(body.Match()).Pipe,
		result,
	); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *EventHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	var result *[]models.Event
	if result, err = dao.EventGet(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Listed(result)
}

func (i *EventHandler) GetPublic(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	var result *[]models.Event
	if result, err = dao.EventGetPublic(c.Ctx(), body); err != nil {
		return
	}
	return c.Listed(result)
}

func (i *EventHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Event)
	if err = dao.EventCollection.UpdateOneAggregate(
		c.Ctx(),
		body.Filter(),
		vmdb.UpdateSet(body),
		result,
		models.EventPipeline(token).Match(body.Match()).Pipe,
	); err != nil {
		return
	}
	result.EditorID = token.ID
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/event/update/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Updated(result)
}

func (i *EventHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.EventDelete(c.Ctx(), body.ID); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
