package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type ParticipationHandler struct {
	vcago.Handler
}

var Participation = &ParticipationHandler{*vcago.NewHandler("participation")}

func (i *ParticipationHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.GET("/user", i.GetByUser, accessCookie)
	group.GET("/event/:id", i.GetByEvent, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.DELETE("/:id", i.DeleteByID, accessCookie)

}

// Create
// @Security CookieAuth
// @Summary Create a Participation
// @Description creates an Participation object.
// @Tags /events/participation
// @Accept json
// @Produce json
// @Param form body models.ParticipationCreate true "Participation Data"
// @Model: vcago.Response
// @Success 201 {object} vcago.ResponseCreated{payload=models.Participation} "Participation successfully created"
// @Failure 400 {object} vcago.BindErrorResponse{} "Bind Error"
// @Failure 400 {object} vcago.ValidationErrorResponse{} "Validation Error"
// @Failure 409 {object} vcago.MongoDuplicatedErrorResponse{} "Duplicated Key"
// @Router /events/participation [post]
func (i *ParticipationHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ParticipationCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Participation)
	if result, err = dao.ParticipationInsert(c.Ctx(), body, token); err != nil {
		return
	}
	dao.ParticipationCreateNotification(c.Ctx(), result)
	return c.Created(result)
}

// Get
// @Security CookieAuth
// @Summary Get a List ofParticipation
// @Tags /events/participation
// @Accept json
// @Produce json
// @Param   q query   models.ParticipationQuery   false  "string collection"  collectionFormat(multi)
// @Model: vcago.Response
// @Success 200 {object} vcago.ResponseListed{payload=[]models.Participation}
// @Router /events/participation [get]
func (i *ParticipationHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ParticipationQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new([]models.Participation)
	if result, err = dao.ParticipationGet(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

// GetByID
// @Security CookieAuth
// @Summary Get a Participation by ID
// @Tags /events/participation
// @Accept json
// @Produce json
// @Param id path string true "Participation ID"
// @Model: vcago.Response
// @Success 200 {object} vcago.ResponseSelected{payload=models.Participation}
// @Failure 404 {object} vcago.MongoNoDocumentErrorResponse{} "No Document with given ID"
// @Router /events/participation/{id} [get]
func (i *ParticipationHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ParticipationParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Participation)
	if result, err = dao.ParticipationGetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *ParticipationHandler) GetByUser(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ParticipationQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new([]models.UserParticipation)
	if result, err = dao.ParticipationUserGet(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *ParticipationHandler) GetByEvent(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new([]models.EventParticipation)
	if result, err = dao.ParticipationEventGet(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

// Update
// @Security CookieAuth
// @Summary Update a Participation
// @Tags /events/participation
// @Accept json
// @Produce json
// @Param form body models.ParticipationUpdate true "Participation Data"
// @Model: vcago.Response
// @Success 200 {object} vcago.ResponseUpdated{payload=models.Participation}
// @Failure 400 {object} vcago.BindErrorResponse{} "Bind Error"
// @Failure 400 {object} vcago.ValidationErrorResponse{} "Validation Error"
// @Failure 404 {object} vcago.MongoNoDocumentErrorResponse{} "No Document with given ID"
// @Router /events/participation [put]
func (i *ParticipationHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ParticipationUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Participation)
	if result, err = dao.ParticipationUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Updated(result)
}

// DeleteByID
// @Security CookieAuth
// @Summary Delete a Participation by ID
// @Tags /events/participation
// @Accept json
// @Produce json
// @Param id path string true "Participation ID"
// @Success 200 {object} vcago.ResponseDeleted{payload=string}
// @Failure 404 {object} vcago.MongoNoDocumentErrorResponse{} "No Document with given ID"
// @Router /events/participation/{id} [delete]
func (i *ParticipationHandler) DeleteByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.ParticipationParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.ParticipationDelete(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
