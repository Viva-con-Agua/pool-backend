package models

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type AvatarCreate struct {
	URL  string `bson:"url" json:"url"`
	Type string `bson:"type" json:"type"`
}

func (i *AvatarCreate) Avatar(userID string) *Avatar {
	return &Avatar{
		ID:       uuid.NewString(),
		URL:      i.URL,
		Type:     i.Type,
		UserID:   userID,
		Modified: vcago.NewModified(),
	}
}

type AvatarUpdate struct {
	ID   string `bson:"_id" json:"id"`
	URL  string `bson:"url" json:"url"`
	Type string `bson:"type" json:"type"`
}
type AvatarParam struct {
	ID string `param:"_id"`
}

type Avatar struct {
	ID       string         `bson:"_id" json:"id"`
	URL      string         `bson:"url" json:"url"`
	Type     string         `bson:"type" json:"type"`
	UserID   string         `bson:"user_id" json:"user_id"`
	Modified vcago.Modified `bson:"modified" json:"modified"`
}

func (i *AvatarUpdate) Filter(token *vcapool.AccessToken) bson.D {
	return bson.D{{Key: "_id", Value: i.ID}, {Key: "user_id", Value: token.ID}}
}

func (i *AvatarParam) Filter(token *vcapool.AccessToken) bson.D {
	return bson.D{{Key: "_id", Value: i.ID}, {Key: "user_id", Value: token.ID}}
}
