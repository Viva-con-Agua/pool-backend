package token

import (
	"pool-user/dao"
	"pool-user/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type AvatarHandler struct {
	vcago.Handler
}

var Avatar = &AvatarHandler{*vcago.NewHandler("avatar")}

func (i *AvatarHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)
}

func (i *AvatarHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AvatarCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := body.Avatar(token.ID)
	if err = dao.AvatarCollection.InsertOne(c.Ctx(), result); err != nil {
		return
	}
	return c.Created(result)
}

func (i *AvatarHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AvatarUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Avatar)
	if err = dao.AvatarCollection.UpdateOne(
		c.Ctx(),
		body.Filter(token),
		vmdb.NewUpdateSet(body),
		result,
	); err != nil {
		return
	}
	return c.Updated(result)

}

func (i *AvatarHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AvatarParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.AvatarCollection.DeleteOne(c.Ctx(), body.Filter(token)); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
