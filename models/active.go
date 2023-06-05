package models

import (
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	Active struct {
		ID       string        `bson:"_id" json:"id"`
		Status   string        `bson:"status" json:"status"`
		Since    int64         `bson:"since" json:"since"`
		UserID   string        `bson:"user_id" json:"user_id"`
		CrewID   string        `bson:"crew_id" json:"crew_id"`
		Modified vmod.Modified `bson:"modified" json:"modified"`
	}
	ActiveUpdate struct {
		Status string `bson:"status" json:"status"`
		Since  int64  `bson:"since" json:"since"`
	}
	ActiveParam struct {
		UserID string `json:"user_id"`
	}
)

var ActiveCollection = "active"

func NewActive(userID string, crewID string) *Active {
	return &Active{
		ID:       uuid.NewString(),
		Status:   "not_requested",
		Since:    time.Now().Unix(),
		UserID:   userID,
		CrewID:   crewID,
		Modified: vmod.NewModified(),
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

func ActiveRequestPermission(token *vcapool.AccessToken) (err error) {
	if token.CrewID == "" {
		return vcago.NewBadRequest("active", "not an crew member")
	}
	return
}

func ActivePermission(token *vcapool.AccessToken) (err error) {
	if !token.Roles.Validate("employee;admin") && !token.PoolRoles.Validate("network;operation") {
		return vcago.NewBadRequest("active", "permission denied")
	}
	return
}

func (i *ActiveParam) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("user_id", i.UserID)
	if !token.Roles.Validate("employee;admin") {
		filter.EqualString("crew_id", token.CrewID)
	}
	return filter.Bson()
}

func (i *Active) ToContent(crew *Crew) *vmod.Content {
	content := &vmod.Content{
		Fields: make(map[string]interface{}),
	}
	content.Fields["Crew"] = crew
	return content
}
