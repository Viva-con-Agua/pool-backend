package models

import (
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
)

type (
	ActivityDatabase struct {
		ID        string        `json:"id" bson:"_id"`
		UserID    string        `json:"user_id" bson:"user_id"`
		Comment   string        `json:"comment" bson:"comment"`
		ModelType string        `json:"model_type" bson:"model_type"`
		ModelID   string        `json:"model_id" bson:"model_id"`
		Status    string        `json:"status" bson:"status"`
		Modified  vmod.Modified `json:"modified" bson:"modified"`
	}
	Activity struct {
		ID        string        `json:"id" bson:"_id"`
		UserID    string        `json:"user_id" bson:"user_id"`
		User      UserDatabase  `json:"user" bson:"user"`
		Comment   string        `json:"comment" bson:"comment"`
		ModelType string        `json:"model_type" bson:"model_type"`
		ModelID   string        `json:"model_id" bson:"model_id"`
		Status    string        `json:"status" bson:"status"`
		Modified  vmod.Modified `json:"modified" bson:"modified"`
	}
)

func NewActivity(userID string, modelType string, modelID string, comment string, status string) *ActivityDatabase {
	return &ActivityDatabase{
		ID:        uuid.NewString(),
		UserID:    userID,
		Comment:   comment,
		ModelType: modelType,
		ModelID:   modelID,
		Status:    status,
		Modified:  vmod.NewModified(),
	}
}

func (i *ActivityDatabase) New(userID string, modelID string) *ActivityDatabase {
	i.ID = uuid.NewString()
	i.UserID = userID
	i.ModelID = modelID
	i.Modified = vmod.NewModified()
	return i
}
