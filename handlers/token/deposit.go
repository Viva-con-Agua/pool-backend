package token

import (
	"net/http"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type DepositResponseSuccess struct {
	vcago.Response
	Model string `json:"model" example:"deposit"`
	Type  string `json:"type" example:"success"`
}

type DepositHandler struct {
	vcago.Handler
}

var Deposit = &DepositHandler{*vcago.NewHandler("deposit")}

func (i *DepositHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.GET("/sync/:id", i.Sync, accessCookie)
}

// Create
// @Security CookieAuth
// @Summary Create a Deposit
// @Description creates an Deposit object.
// @Tags Deposit
// @Accept json
// @Produce json
// @Param form body models.DepositCreate true "Deposit Data"
// @Model: vcago.Response
// @Success 201 {object} vcago.ResponseCreated{payload=models.Deposit} "Successfully Created Deposit"
// @Failure 400 {object} vcago.BindErrorResponse{} "Bind Error"
// @Failure 400 {object} vcago.ValidationErrorResponse{} "Validation Error"
// @Failure 409 {object} vcago.MongoDuplicatedErrorResponse{} "Duplicated Key"
// @Router /finances/deposit [post]
func (i *DepositHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.DepositCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Deposit)
	if result, err = dao.DepositInsert(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Created(result)
}

// Get
// @Security CookieAuth
// @Summary Get a List of Deposit
// @Tags Deposit
// @Accept json
// @Produce json
// @Param   q query   models.DepositQuery   false  "string collection"  collectionFormat(multi)
// @Model: vcago.Response
// @Success 200 {object} vcago.ResponseListed{payload=[]models.Deposit}
// @Failure 400 {object} vcago.Response{}
// @Router /finances/deposit [get]
func (i *DepositHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.DepositQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new([]models.Deposit)
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if result, err = dao.DepositGet(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Listed(result, int64(len(*result)))
}

// GetByID
// @Security CookieAuth
// @Summary Get a Deposit by ID
// @Tags Deposit
// @Accept json
// @Produce json
// @Param id path string true "Deposit ID"
// @Model: vcago.Response
// @Success 200 {object} vcago.ResponseSelected{payload=models.Deposit}
// @Router /finances/deposit/{id} [get]
func (i *DepositHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.DepositParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Deposit)
	if result, err = dao.DepositGetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

// Update
// @Security CookieAuth
// @Summary Update a Deposit
// @Tags Deposit
// @Accept json
// @Produce json
// @Param form body models.DepositUpdate true "Deposit Data"
// @Model: vcago.Response
// @Success 200 {object} vcago.ResponseUpdated{payload=models.Deposit}
// @Router /finances/deposit [put]
func (i *DepositHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.DepositUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Deposit)
	if result, err = dao.DepositUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Updated(result)
}

// Sync
// @Security CookieAuth
// @Summary Sync Deposit by ID
// @Tags Deposit
// @Accept json
// @Produce json
// @Param id path string true "Deposit ID"
// @Model: vcago.Response
// @Success 200 {object} vcago.ResponseSynced{}
// @Router /finances/deposit/sync/{id} [get]
func (i *DepositHandler) Sync(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.DepositParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = body.DepositSyncPermission(token); err != nil {
		return
	}
	if _, err = dao.DepositSync(c.Ctx(), body, token); err != nil {
		return
	}
	return c.SuccessResponse(http.StatusOK, "successfully_synced", "event", nil)
}

// DeleteByID
// @Security CookieAuth
// @Summary Delete a Deposit by ID
// @Tags Deposit
// @Accept json
// @Produce json
// @Param id path string true "Deposit ID"
// @Success 200 {object} vcago.ResponseDeleted{payload=string}
// @Router /fincances/deposit/{id} [delete]
func (i *DepositHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.DepositParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	return c.Deleted(body.ID)

}
