package token

import (
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/labstack/echo/v4"
)

type TakingHandler struct {
	vcago.Handler
}

var Taking = &TakingHandler{*vcago.NewHandler("taking")}

func (i *TakingHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.DELETE("/:id", i.DeleteByID, accessCookie)

}

// Create
// @Security CookieAuth
// @Summary Create a Taking
// @Description creates an Taking object.
// @Tags /finances/taking
// @Accept json
// @Produce json
// @Param form body models.TakingCreate true "Taking Data"
// @Model: vcago.Response
// @Success 201 {object} vcago.ResponseCreated{payload=models.Taking} "Taking succsessfully created"
// @Failure 400 {object} vcago.BindErrorResponse{} "Bind Error"
// @Failure 400 {object} vcago.ValidationErrorResponse{} "Validation Error"
// @Failure 409 {object} vcago.MongoDuplicatedErrorResponse{} "Duplicated Key"
// @Router  /finances/taking [post]
func (i *TakingHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.TakingCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Taking)
	if result, err = dao.TakingInsert(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Created(result)
}

// Update
// @Security CookieAuth
// @Summary Update a Taking
// @Tags /finances/taking
// @Accept json
// @Produce json
// @Param form body models.TakingUpdate true "Taking Data"
// @Model: vcago.Response
// @Success 200 {object} vcago.ResponseUpdated{payload=models.Taking}
// @Failure 400 {object} vcago.BindErrorResponse{} "Bind Error"
// @Failure 400 {object} vcago.ValidationErrorResponse{} "Validation Error"
// @Failure 404 {object} vcago.MongoNoDocumentErrorResponse{} "No Document with given ID"
// @Router  /finances/taking [put]
func (i TakingHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.TakingUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Taking)
	if result, err = dao.TakingUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Updated(result)
}

// Get
// @Security CookieAuth
// @Summary Get a List ofTaking
// @Tags /finances/taking
// @Accept json
// @Produce json
// @Param   q query   models.TakingQuery   false  "string collection"  collectionFormat(multi)
// @Model: vcago.Response
// @Success 200 {object} vcago.ResponseListed{payload=[]models.Taking}
// @Router  /finances/taking [get]
func (i TakingHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.TakingQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	var result []models.Taking
	var listSize int64
	if result, listSize, err = dao.TakingGet(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Listed(result, listSize)
}

// GetByID
// @Security CookieAuth
// @Summary Get a Taking by ID
// @Tags /finances/taking
// @Accept json
// @Produce json
// @Param id path string true "Taking ID"
// @Model: vcago.Response
// @Success 200 {object} vcago.ResponseSelected{payload=models.Taking}
// @Failure 404 {object} vcago.MongoNoDocumentErrorResponse{} "No Document with given ID"
// @Router  /finances/taking/{id} [get]
func (i TakingHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(vmod.IDParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Taking)
	if result, err = dao.TakingGetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

// DeleteByID
// @Security CookieAuth
// @Summary Delete a Taking by ID
// @Tags /finances/taking
// @Accept json
// @Produce json
// @Param id path string true "Taking ID"
// @Success 200 {object} vcago.ResponseDeleted{payload=string}
// @Failure 404 {object} vcago.MongoNoDocumentErrorResponse{} "No Document with given ID"
// @Router  /finances/taking/{id} [delete]
func (i TakingHandler) DeleteByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(vmod.IDParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.TakingDeletetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
