package models

import (
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	Active struct {
		ID       string         `bson:"_id" json:"id"`
		Status   string         `bson:"status" json:"status"`
		Since    int64          `bson:"since" json:"since"`
		UserID   string         `bson:"user_id" json:"user_id"`
		CrewID   string         `bson:"crew_id" json:"crew_id"`
		Modified vcago.Modified `bson:"modified" json:"modified"`
	}
	ActiveUpdate struct {
		Status string `bson:"status" json:"status"`
		Since  int64  `bson:"since" json:"since"`
	}
	ActiveParam struct {
		UserID string `json:"user_id"`
	}
)

func NewActive(userID string, crewID string) *Active {
	return &Active{
		ID:       uuid.NewString(),
		Status:   "not_requested",
		Since:    time.Now().Unix(),
		UserID:   userID,
		CrewID:   crewID,
		Modified: vcago.NewModified(),
	}
}

func ActiveConfirm() *ActiveUpdate {
	return &ActiveUpdate{
		Status: "confirmed",
		Since:  time.Now().Unix(),
	}
}

func ActiveReject() *ActiveUpdate {
	return &ActiveUpdate{
		Status: "rejected",
		Since:  time.Now().Unix(),
	}
}

func ActiveWithdraw() *ActiveUpdate {
	return &ActiveUpdate{
		Status: "withdrawn",
		Since:  time.Now().Unix(),
	}
}

func ActiveRequest() *ActiveUpdate {
	return &ActiveUpdate{
		Status: "requested",
		Since:  time.Now().Unix(),
	}
}

func (i *Active) IsRequested() bool {
	if i.Status == "requested" {
		return true
	}
	return false
}
func (i *Active) IsConfirmed() bool {
	if i.Status == "confirmed" {
		return true
	}
	return false
}
func (i *Active) IsWithdrawn() bool {
	if i.Status == "withdrawn" {
		return true
	}
	return false
}
func (i *Active) IsRejected() bool {
	if i.Status == "rejected" {
		return true
	}
	return false
}

func ActiveRequestPermission(token *vcapool.AccessToken) (err error) {
	if token.CrewID == "" {
		return vcago.NewBadRequest("active", "not an crew member")
	}
	return
}

func ActivePermission(token *vcapool.AccessToken) (err error) {
	if !token.Roles.Validate("employee") && !token.PoolRoles.Validate("network;operation") {
		return vcago.NewBadRequest("active", "permission denied")
	}
	return
}

func (i *ActiveParam) Filter(token *vcapool.AccessToken) bson.D {
	if token.Roles.Validate("employee") {
		return bson.D{{Key: "user_id", Value: i.UserID}}

	}
	return bson.D{{Key: "user_id", Value: i.UserID}, {Key: "crew_id", Value: token.CrewID}}
}
