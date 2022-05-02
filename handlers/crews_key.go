package handlers

import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type CrewAPIKeyHandler struct {
	vcago.Handler
}

func NewCrewAPIKeyHandler() *CrewAPIKeyHandler {
	handler := vcago.NewHandler("crew")
	return &CrewAPIKeyHandler{
		*handler,
	}
}

func (i *CrewAPIKeyHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(dao.Crew)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	if err = body.Create(c.Ctx()); err != nil {
		return
	}
	return c.Created(body)
}
