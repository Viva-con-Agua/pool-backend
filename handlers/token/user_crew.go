package token

import (
	"log"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

type UserCrewHandler struct {
	vcago.Handler
}

var UserCrew = &UserCrewHandler{*vcago.NewHandler("user_crew")}

func (i *UserCrewHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.POST("/create", i.UsersCreate, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.PUT("/update", i.UsersUpdate, accessCookie)
	group.DELETE("", i.Delete, accessCookie)
	group.DELETE("/:id", i.UsersDelete, accessCookie)
}

func (i *UserCrewHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserCrewCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.UserCrew)
	if result, err = dao.UserCrewInsert(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/profile/crew/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Created(result)
}

func (i *UserCrewHandler) UsersCreate(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UsersCrewCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.UserCrew)
	if result, err = dao.UsersUserCrewInsert(c.Ctx(), body, token); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/profile/crew/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Created(result)
}

func (i *UserCrewHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserCrewUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.UserCrew)
	if result, err = dao.UserCrewUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Updated(result)
}

func (i *UserCrewHandler) UsersUpdate(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserCrewUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.UserCrew)
	if result, err = dao.UsersCrewUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Updated(result)
}

func (i *UserCrewHandler) UsersDelete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = models.UsersEditPermission(token); err != nil {
		return
	}
	if err = dao.UserCrewDelete(c.Ctx(), body.ID); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(&models.UserCrewUpdate{UserID: token.ID}, "/v1/pool/profile/crew/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Deleted(token.ID)

}

func (i *UserCrewHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.UserCrewDelete(c.Ctx(), token.ID); err != nil {
		return
	}
	go func() {
		if err = dao.IDjango.Post(&models.UserCrewUpdate{UserID: token.ID}, "/v1/pool/profile/crew/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Deleted(token.ID)

}
