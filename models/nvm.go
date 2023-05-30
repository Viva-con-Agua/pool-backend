package models

import (
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
)

var NVMCollection = "nvm"

type (
	NVM struct {
		ID       string        `bson:"_id" json:"id"`
		Status   string        `bson:"status" json:"status"`
		Since    int64         `bson:"since" json:"since"`
		UserID   string        `bson:"user_id" json:"user_id"`
		Modified vmod.Modified `bson:"modified" json:"modified"`
	}
	NVMUpdate struct {
		Status string `bson:"status" json:"status"`
		Since  int64  `bson:"since" json:"since"`
	}
	NVMParam struct {
		UserID string `json:"user_id"`
	}
	NVMIDParam struct {
		ID string `param:"id"`
	}
)

func NewNVM(userID string) *NVM {
	return &NVM{
		ID:       uuid.NewString(),
		Status:   "not_requested",
		Since:    time.Now().Unix(),
		UserID:   userID,
		Modified: vmod.NewModified(),
	}
}

func NewNVMRejected() *NVMUpdate {
	return &NVMUpdate{
		Status: "rejected",
		Since:  time.Now().Unix(),
	}
}

func NVMConfirmedPermission(token *vcapool.AccessToken) (err error) {
	if token.ActiveState != "confirmed" {
		return vcago.NewBadRequest("user_nvm", "active required")
	}
	if token.AddressID == "" {
		return vcago.NewBadRequest("user_nvm", "address required")
	}
	if token.Birthdate == 0 {
		return vcago.NewBadRequest("user_nvm", "birthdate required")
	}
	return
}

func NVMPermission(token *vcapool.AccessToken) (err error) {
	if !token.Roles.Validate("employee;admin") {
		return vcago.NewPermissionDenied("nvm")
	}
	return
}

func NVMConfirm() *NVMUpdate {
	return &NVMUpdate{
		Status: "confirmed",
		Since:  time.Now().Unix(),
	}
}

func NVMReject() *NVMUpdate {
	return &NVMUpdate{
		Status: "rejected",
		Since:  time.Now().Unix(),
	}
}

func NVMWithdraw() *NVMUpdate {
	return &NVMUpdate{
		Status: "withdrawn",
		Since:  time.Now().Unix(),
	}
}

func (i *NVM) IsConfirmed() bool {
	if i.Status == "confirmed" {
		return true
	}
	return false
}
func (i *NVM) IsWithdrawn() bool {
	if i.Status == "withdrawn" {
		return true
	}
	return false
}
func (i *NVM) IsRejected() bool {
	if i.Status == "rejected" {
		return true
	}
	return false
}
