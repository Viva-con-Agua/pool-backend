package token

import (
	"log"
	"net/http"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type EventHandler struct {
	vcago.Handler
}

var Event = &EventHandler{*vcago.NewHandler("event")}

func (i *EventHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.POST("", i.Create, accessCookie)
	group.PUT("", i.Update, accessCookie)
	group.DELETE("/:id", i.Delete, accessCookie)
	group.GET("", i.Get, accessCookie)
	group.GET("/:id", i.GetByID, accessCookie)
	group.GET("/sync/:id", i.Sync, accessCookie)
	group.GET("/public", i.GetPublic)
	group.GET("/email", i.GetEmailEvents, accessCookie)
	group.GET("/user", i.GetByEventAsp, accessCookie)
	group.GET("/view/:id", i.GetViewByID)
	group.GET("/details/:id", i.GetPrivateDetails, accessCookie)

}

func (i *EventHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventCreate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Event)
	if result, err = dao.EventInsert(c.Ctx(), body, token); err != nil {
		return
	}
	result.EditorID = token.ID
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/event/create/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Created(result)
}

func (i *EventHandler) Get(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}

	result := new([]models.ListEvent)
	var listSize int64
	if result, listSize, err = dao.EventGet(body, token); err != nil {
		return
	}
	return c.Listed(result, listSize)
}

func (i *EventHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Event)
	if result, err = dao.EventGetByID(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *EventHandler) GetViewByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new(models.EventPublic)
	if result, err = dao.EventViewGetByID(c.Ctx(), body); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *EventHandler) GetPrivateDetails(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.EventDetails)
	if result, err = dao.ParticipationAspGet(c.Ctx(), &models.ParticipationQuery{EventID: []string{body.ID}}, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *EventHandler) GetPublic(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	result := new([]models.EventPublic)
	var listSize int64
	if result, listSize, err = dao.EventGetPublic(body); err != nil {
		return
	}
	return c.Listed(result, listSize)
}

func (i *EventHandler) GetByEventAsp(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new([]models.ListDetailsEvent)
	if result, err = dao.EventGetAps(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *EventHandler) GetEmailEvents(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventQuery)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new([]models.EventPublic)
	if result, err = dao.EventsGetReceiverEvents(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Selected(result)
}

func (i *EventHandler) Update(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventUpdate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.Event)
	if result, err = dao.EventUpdate(c.Ctx(), body, token); err != nil {
		return
	}
	result.EditorID = token.ID
	go func() {
		if err = dao.IDjango.Post(result, "/v1/pool/event/update/"); err != nil {
			log.Print(err)
		}
	}()
	return c.Updated(result)
}

func (i *EventHandler) Sync(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = body.EventSyncPermission(token); err != nil {
		return
	}
	if _, err = dao.EventSync(c.Ctx(), body, token); err != nil {
		return
	}
	return c.SuccessResponse(http.StatusOK, "successfully_synced", "event", nil)
}

func (i *EventHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.EventParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = dao.EventDelete(c.Ctx(), body, token); err != nil {
		return
	}
	return c.Deleted(body)
}
