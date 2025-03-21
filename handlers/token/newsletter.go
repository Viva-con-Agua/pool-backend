package token

import (
	"log"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
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

// Create
// @Security CookieAuth
// @Summary Create a Newsletter
// @Description creates an Newsletter object.
// @Tags /users/newsletter
// @Accept json
// @Produce json
// @Param form body models.Newsletter true "Newsletter Data"
// @Model: vcago.Response
// @Success 201 {object} vcago.ResponseCreated{payload=models.Newsletter} "Newsletter successfully created"
// @Failure 400 {object} vcago.BindErrorResponse{} "Bind Error"
// @Failure 400 {object} vcago.ValidationErrorResponse{} "Validation Error"
// @Failure 409 {object} vcago.MongoDuplicatedErrorResponse{} "Duplicated Key"
// @Router /users/newsletter [post]
func (i *NewsletterHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.NewsletterCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Newsletter)
	if result, err = dao.NewsletterCreate(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/newsletter/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Created(result)
}

// DeleteByID
// @Security CookieAuth
// @Summary Delete a Newsletter by ID
// @Tags /users/newsletter
// @Accept json
// @Produce json
// @Param id path string true "Newsletter ID"
// @Success 200 {object} vcago.ResponseDeleted{payload=string} "Newsletter successfully deleted"
// @Failure 404 {object} vcago.MongoNoDocumentErrorResponse{} "No Document with given ID"
// @Router /users/newsletter/{id} [delete]
func (i *NewsletterHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.NewsletterParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Newsletter)
	if result, err = dao.NewsletterDelete(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/newsletter/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Deleted(body.ID)
}
