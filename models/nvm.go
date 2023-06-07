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

func NVMConfirmedPermission(token *vcapool.AccessToken) (err error) {
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

func NVMPermission(token *vcapool.AccessToken) (err error) {
	if !token.Roles.Validate("employee;admin") {
		return vcago.NewPermissionDenied(NVMCollection)
	}
	return
}

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
