package models

import (
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

var UserCrewCollection = "user_crew"

type (
	UserCrewCreate struct {
		CrewID string `json:"crew_id"`
	}
	UserCrew struct {
		ID        string        `bson:"_id" json:"id"`
		UserID    string        `bson:"user_id" json:"user_id"`
		Name      string        `bson:"name" json:"name"`
		Email     string        `bson:"email" json:"email"`
		Roles     []vmod.Role   `bson:"roles" json:"roles"`
		CrewID    string        `bson:"crew_id" json:"crew_id"`
		MailboxID string        `bson:"mailbox_id" json:"mailbox_id"`
		Modified  vmod.Modified `bson:"modified" json:"modified"`
	}
	UserCrewUpdate struct {
		ID     string `bson:"_id" json:"id"`
		UserID string `bson:"user_id" json:"user_id"`
		Name   string `bson:"name" json:"name"`
		Email  string `bson:"email" json:"email"`
		CrewID string `bson:"crew_id" json:"crew_id"`
	}
	UserCrewParam struct {
		ID string `param:"id"`
	}
	UserCrewImport struct {
		DropsID string `json:"drops_id"`
		CrewID  string `json:"crew_id"`
	}
)

func NewUserCrew(userID string, crewID string, name string, email string, mailboxID string) *UserCrew {
	return &UserCrew{
		ID:        uuid.NewString(),
		UserID:    userID,
		Name:      name,
		Email:     email,
		CrewID:    crewID,
		MailboxID: mailboxID,
		Modified:  vmod.NewModified(),
	}
}

func (i *UserCrewCreate) CrewFilter() bson.D {
	match := vmdb.NewFilter()
	match.EqualString("_id", i.CrewID)
	return bson.D(*match)
}

func (i *UserCrewUpdate) Filter(token *vcapool.AccessToken) bson.D {
	return bson.D{{Key: "_id", Value: i.ID}, {Key: "user_id", Value: token.ID}}
}
