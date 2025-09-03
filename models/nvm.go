package models

import (
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmod"
)

var NVMCollection = "nvm"

type (
	NVM struct {
		Status   string        `bson:"status" json:"status"`
		Since    int64         `bson:"since" json:"since"`
		UserID   string        `bson:"user_id" json:"user_id"`
		Modified vmod.Modified `bson:"modified" json:"modified"`
	}
	NVMUpdate struct {
		Status  string `bson:"nvm.status" json:"status"`
		Since   int64  `bson:"nvm.since" json:"since"`
		Updated int64  `bson:"nvm.modified.updated"`
	}
	NVMParam struct {
		UserID string `json:"user_id"`
	}
	NVMIDParam struct {
		ID string `param:"id"`
	}
)

func NVMConfirmedPermission(token *AccessToken) (err error) {
	if token.ActiveState != "confirmed" {
		return vcago.NewBadRequest(NVMCollection, "active required")
	}
	if token.AddressID == "" {
		return vcago.NewBadRequest(NVMCollection, "address required")
	}
	if token.Birthdate == 0 {
		return vcago.NewBadRequest(NVMCollection, "birthdate required")
	}
	return
}

func NVMPermission(token *AccessToken) (err error) {
	if !token.Roles.Validate("admin;employee;pool_employee") {
		return vcago.NewPermissionDenied(NVMCollection)
	}
	return
}

func NewNVM(userID string) *NVM {
	return &NVM{
		Status:   "not_requested",
		Since:    time.Now().Unix(),
		UserID:   userID,
		Modified: vmod.NewModified(),
	}
}

func NewNVMRejected() *NVMUpdate {
	return &NVMUpdate{
		Status:  "rejected",
		Since:   time.Now().Unix(),
		Updated: time.Now().Unix(),
	}
}

func NVMConfirm() *NVMUpdate {
	return &NVMUpdate{
		Status:  "confirmed",
		Since:   time.Now().Unix(),
		Updated: time.Now().Unix(),
	}
}

func NVMReject() *NVMUpdate {
	return &NVMUpdate{
		Status:  "rejected",
		Since:   time.Now().Unix(),
		Updated: time.Now().Unix(),
	}
}

func NVMWithdraw() *NVMUpdate {
	return &NVMUpdate{
		Status:  "withdrawn",
		Since:   time.Now().Unix(),
		Updated: time.Now().Unix(),
	}
}
