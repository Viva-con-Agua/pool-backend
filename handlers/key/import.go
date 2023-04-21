package key

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type ImportHandler struct {
	vcago.Handler
}

var Import = &ImportHandler{*vcago.NewHandler("crew")}

func (i *ImportHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("/crew", i.ImportCrew, vcago.KeyAuthMiddleware())
	group.POST("/profile", i.ImportProfile, vcago.KeyAuthMiddleware())
	group.POST("/address", i.ImportAddress, vcago.KeyAuthMiddleware())
	group.POST("/usercrew", i.ImportUserCrew, vcago.KeyAuthMiddleware())
}

func (i *ImportHandler) ImportCrew(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.Crew)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	var result *models.Crew
	if result, err = dao.CrewInsert(c.Ctx(), body); err != nil {
		return
	}
	return c.Created(result)
}

func (i *ImportHandler) ImportProfile(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ProfileImport)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	var result *models.Profile
	if result, err = dao.ProfileImport(c.Ctx(), body); err != nil {
		return
	}
	return c.Created(result)
}

func (i *ImportHandler) ImportAddress(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AddressImport)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	var result *models.Address
	if result, err = dao.AddressImport(c.Ctx(), body); err != nil {
		return
	}
	return c.Created(result)
}

func (i *ImportHandler) ImportUserCrew(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.UserCrewImport)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	var result *models.UserCrew
	if result, err = dao.UserCrewImport(c.Ctx(), body); err != nil {
		return
	}
	return c.Created(result)
}
