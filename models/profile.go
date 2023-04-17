package models

import (
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

var ProfilesCollection = "profiles"

type (
	ProfileCreate struct {
		Gender     string `bson:"gender" json:"gender"`
		Phone      string `bson:"phone" json:"phone"`
		Mattermost string `bson:"mattermost_username" json:"mattermost_username"`
		Birthdate  int64  `bson:"birthdate" json:"birthdate"`
	}
	ProfileUpdate struct {
		ID         string `bson:"_id" json:"id"`
		Gender     string `bson:"gender" json:"gender"`
		Phone      string `bson:"phone" json:"phone"`
		Mattermost string `bson:"mattermost_username" json:"mattermost_username"`
		Birthdate  int64  `bson:"birthdate" json:"birthdate"`
	}
	Profile struct {
		ID         string        `bson:"_id" json:"id"`
		Gender     string        `bson:"gender" json:"gender"`
		Phone      string        `bson:"phone" json:"phone"`
		Mattermost string        `bson:"mattermost_username" json:"mattermost_username"`
		Birthdate  int64         `bson:"birthdate" json:"birthdate"`
		UserID     string        `bson:"user_id" json:"user_id"`
		Modified   vmod.Modified `bson:"modified" json:"modified"`
	}
)

func (i *ProfileCreate) Profile(userID string) *Profile {
	return &Profile{
		ID:         uuid.NewString(),
		Gender:     i.Gender,
		Phone:      i.Phone,
		Mattermost: i.Mattermost,
		Birthdate:  i.Birthdate,
		UserID:     userID,
		Modified:   vmod.NewModified(),
	}
}

func (i *ProfileUpdate) Filter(token *vcapool.AccessToken) bson.D {
	return bson.D{{Key: "_id", Value: i.ID}, {Key: "user_id", Value: token.ID}}
}
