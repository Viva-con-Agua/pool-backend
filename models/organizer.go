package models

import (
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	OrganizerCreate struct {
		Name string `json:"name" bson:"name" validate:"required"`
	}
	Organizer struct {
		ID       string        `json:"id" bson:"_id"`
		Name     string        `json:"name" bson:"name"`
		Modified vmod.Modified `json:"modified" bson:"modified"`
	}
	OrganizerUpdate struct {
		ID   string `json:"id" bson:"_id" validate:"required"`
		Name string `json:"name" bson:"name" validate:"required"`
	}
	OrganizerParam struct {
		ID string `param:"id"`
	}
	OrganizerQuery struct {
		ID          string `query:"id" qs:"id"`
		Name        string `query:"name" qs:"name"`
		UpdatedTo   string `query:"updated_to" qs:"updated_to"`
		UpdatedFrom string `query:"updated_from" qs:"updated_from"`
		CreatedTo   string `query:"created_to" qs:"created_to"`
		CreatedFrom string `query:"created_from" qs:"created_from"`
	}
)

func (i *OrganizerCreate) Organizer() *Organizer {
	return &Organizer{
		ID:       uuid.NewString(),
		Name:     i.Name,
		Modified: vmod.NewModified(),
	}
}

func (i *OrganizerParam) Filter() bson.D {
	return bson.D{{Key: "_id", Value: i.ID}}
}

func (i *OrganizerUpdate) Filter() bson.D {
	return bson.D{{Key: "_id", Value: i.ID}}
}

func (i *OrganizerQuery) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	filter.LikeString("name", i.Name)
	filter.GteInt64("modified.updated", i.UpdatedFrom)
	filter.GteInt64("modified.created", i.CreatedFrom)
	filter.LteInt64("modified.updated", i.UpdatedTo)
	filter.LteInt64("modified.created", i.CreatedTo)
	return bson.D(*filter)
}
