package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type SourceHandler struct {
	vcago.Handler
}

var Source = &SourceHandler{*vcago.NewHandler("source")}

func (i *SourceHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)
}

func (i *SourceHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.SourceCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := body.Source()
	if err = dao.SourceCollection.InsertOne(c.Ctx(), result); err != nil {
		return
	}
	return c.Created(result)

}

func (i *SourceHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.SourceUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Source)
	if err = dao.SourceCollection.UpdateOne(
		c.Ctx(),
		body.Filter(),
		vmdb.UpdateSet(body),
		result,
	); err != nil {
		return
	}
	return c.Updated(result)
}

func (i *SourceHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.SourceQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new([]models.Source)
	if err = dao.SourceCollection.Find(c.Ctx(), body.Filter(), result); err != nil {
		return
	}
	return c.Listed(result)
}

func (i SourceHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.SourceParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new(models.Source)
	if err = dao.SourceCollection.FindOne(c.Ctx(), body.Filter(), result); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *SourceHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.SourceParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	if err = dao.SourceCollection.DeleteOne(c.Ctx(), body.Filter()); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
