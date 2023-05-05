package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type ParticipationHandler struct {
	vcago.Handler
}

var Participation = &ParticipationHandler{*vcago.NewHandler("participation")}

func (i *ParticipationHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.PUT("/status", i.Status, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)

}

func (i *ParticipationHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ParticipationCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	database := body.ParticipationDatabase(token)
	if err = dao.ParticipationCollection.InsertOne(c.Ctx(), database); err != nil {
		return
	}
	result := new(models.Participation)
	if err = dao.ParticipationCollection.AggregateOne(
		c.Ctx(),
		models.ParticipationPipeline().Match(database.Match()).Pipe,
		result,
	); err != nil {
		return
	}
	return c.Created(result)
}

func (i *ParticipationHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ParticipationParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Participation)
	if err = dao.ParticipationCollection.AggregateOne(
		c.Ctx(),
		models.ParticipationPipeline().Match(body.Match()).Pipe,
		result,
	); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *ParticipationHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ParticipationUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Participation)
	if err = dao.ParticipationCollection.UpdateOneAggregate(
		c.Ctx(),
		body.Filter(),
		vmdb.UpdateSet(body),
		result,
		models.ParticipationPipeline().Match(body.Match()).Pipe,
	); err != nil {
		return
	}
	return c.Updated(result)
}

func (i *ParticipationHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ParticipationQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new([]models.Participation)
	if err = dao.ParticipationCollection.Aggregate(
		c.Ctx(),
		models.ParticipationPipeline().Match(body.Match()).Pipe,
		result,
	); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *ParticipationHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ParticipationParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.ParticipationCollection.DeleteOne(c.Ctx(), body.Filter()); err != nil {
		return
	}
	return c.Deleted(body.ID)
}

func (i *ParticipationHandler) Status(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ParticipationStateRequest)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Participation)
	if err = dao.ParticipationCollection.UpdateOneAggregate(
		c.Ctx(),
		body.Permission(token),
		vmdb.UpdateSet(body),
		result,
		models.ParticipationPipeline().Match(body.Match()).Pipe,
	); err != nil {
		return
	}
	return c.Selected(result)
}
