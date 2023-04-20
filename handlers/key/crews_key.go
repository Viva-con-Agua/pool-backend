package key

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type CrewHandler struct {
	vcago.Handler
}

var Crew = &CrewHandler{*vcago.NewHandler("crew")}

func (i *CrewHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, vcago.KeyAuthMiddleware())
}

func (i *CrewHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.Crew)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	var result *models.Crew
	if result, err = dao.CrewInsert(c.Ctx(), body); err != nil {
		return c.ErrorResponse(err)
	}
	return c.Created(result)
}
