package token

import (
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
	group.GET("", i.Get)
	group.GET("/:id", i.GetByID)
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
	if err = dao.EventCollection.InsertOne(c.Ctx(), database); err != nil {
		return
	}
	result := new(models.Event)
	if err = dao.EventCollection.AggregateOne(
		c.Ctx(),
		models.EventPipeline().Match(database.Match()).Pipe,
		result,
	); err != nil {
		return
	}
	return c.Created(result)
}

func (i *EventHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new(models.Event)
	if err = dao.EventCollection.AggregateOne(
		c.Ctx(),
		models.EventPipeline().Match(body.Match()).Pipe,
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
	result := new([]models.Event)
	if err = dao.EventCollection.Aggregate(
		c.Ctx(),
		models.EventPipeline().Match(body.Match()).Pipe,
		result,
	); err != nil {
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
	result := new(models.Event)
	if err = dao.EventCollection.UpdateOneAggregate(
		c.Ctx(),
		body.Filter(),
		vmdb.UpdateSet(body),
		result,
		models.EventPipeline().Match(body.Match()).Pipe,
	); err != nil {
		return
	}
	return c.Updated(result)
}

func (i *EventHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	if err = dao.EventCollection.DeleteOne(c.Ctx(), body.Filter()); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
