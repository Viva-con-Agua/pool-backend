package token

import (
	"net/http"
	"pool-user/dao"
	"pool-user/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type NVMHandler struct {
	vcago.Handler
}

var NVM = &NVMHandler{*vcago.NewHandler("nvm")}

func (i *NVMHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.GET("/confirm", i.Confirm, accessCookie)
	group.POST("/reject", i.Reject, accessCookie)
	group.GET("/withdraw", i.Withdraw, accessCookie)
}

func (i *NVMHandler) Confirm(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	if err = models.NVMConfirmedPermission(token); err != nil {
		return
	}
	result := new(models.NVM)
	if err = dao.NVMCollection.UpdateOne(
		c.Ctx(),
		bson.D{{Key: "user_id", Value: token.ID}},
		vmdb.NewUpdateSet(models.NVMConfirm()),
		result,
	); err != nil {
		return
	}
	return c.SuccessResponse(http.StatusOK, "successfully_confirmed", "nvm", result)
}

func (i *NVMHandler) Reject(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.NVMParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	//get requested user from token
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	//check if requested user has the network or operation permission
	if err = models.NVMRejectPermission(token); err != nil {
		return
	}
	result := new(models.NVM)
	if err = dao.NVMCollection.UpdateOne(
		c.Ctx(),
		bson.D{{Key: "user_id", Value: body.UserID}},
		vmdb.NewUpdateSet(models.NVMReject()),
		result,
	); err != nil {
		return
	}
	return c.SuccessResponse(http.StatusOK, "successfully_rejected", "nvm", result)
}

func (i *NVMHandler) Withdraw(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(vcapool.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := new(models.NVM)
	if err = dao.NVMCollection.UpdateOne(
		c.Ctx(),
		bson.D{{Key: "user_id", Value: token.ID}},
		vmdb.NewUpdateSet(models.NVMWithdraw()),
		result,
	); err != nil {
		return
	}
	return c.SuccessResponse(http.StatusOK, "successfully_withdrawed", "nvm", result)
}
