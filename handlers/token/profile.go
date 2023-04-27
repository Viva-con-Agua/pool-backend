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

type ProfileHandler struct {
	vcago.Handler
}

var Profile = &ProfileHandler{*vcago.NewHandler("profile")}

func (i *ProfileHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.PUT("", i.Update, accessCookie)
}

func (i *ProfileHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ProfileCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := body.Profile(token.ID)
	if err = dao.ProfilesCollection.InsertOne(c.Ctx(), result); err != nil {
		return
	}
	// TODO: ADD PROFLE CREATE
	return c.Created(result)
}

func (i *ProfileHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ProfileUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Profile)

	if err = dao.ProfilesCollection.UpdateOne(
		c.Ctx(),
		body.Filter(token),
		vmdb.UpdateSet(body),
		result,
	); err != nil {
		return
	}
	if body.Birthdate == 0 {
		var nvm *models.NVM
		if nvm, err = dao.NVMWithdraw(c.Ctx(), token); err != nil {
			return
		}
		if err = dao.IDjango.Post(nvm, "/v1/pool/profile/nvm/"); err != nil {
			log.Print(err)
		}
	}
	if err = dao.IDjango.Post(result, "/v1/pool/profile/"); err != nil {
		log.Print(err)
	}
	return c.Updated(result)
}
