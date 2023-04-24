package token

import (
	"log"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type NewsletterHandler struct {
	vcago.Handler
}

var Newsletter = &NewsletterHandler{*vcago.NewHandler("newsletter")}

func (i *NewsletterHandler) Routes(group *echo.Group) {
	group.POST("", i.Create, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)

}

func (i *NewsletterHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.NewsletterCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	var result *models.Newsletter
	if result, err = dao.NewsletterCreate(c.Ctx(), body, token); err != nil {
		return
	}
	if err = dao.IDjango.Post(result, "/v1/pool/newsletter/"); err != nil {
		log.Print(err)
	}
	return c.Created(result)
}

func (i *NewsletterHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.NewsletterParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	var result *models.Newsletter
	if result, err = dao.NewsletterDelete(c.Ctx(), body, token); err != nil {
		return
	}
	if err = dao.IDjango.Post(result, "/v1/pool/newsletter/"); err != nil {
		log.Print(err)
	}
	return c.Deleted(body.ID)
}
