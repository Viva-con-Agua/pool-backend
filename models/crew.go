package models

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	CrewCreate struct {
		Name   string `json:"name" bson:"name"`
		Email  string `json:"email" bson:"email"`
		Cities []City `json:"cities" bson:"cities"`
	}
	CrewUpdate struct {
		ID     string `json:"id,omitempty" bson:"_id"`
		Name   string `json:"name" bson:"name"`
		Email  string `json:"email" bson:"email"`
		Cities []City `json:"cities" bson:"cities"`
	}
	Crew struct {
		ID       string         `json:"id,omitempty" bson:"_id"`
		Name     string         `json:"name" bson:"name"`
		Email    string         `json:"email" bson:"email"`
		Cities   []City         `json:"cities" bson:"cities"`
		Modified vcago.Modified `json:"modified" bson:"modified"`
	}
	City struct {
		City        string         `json:"city" bson:"city"`
		Country     string         `json:"country" bson:"country"`
		CountryCode string         `json:"country_code" bson:"country_code"`
		PlaceID     string         `json:"place_id" bson:"place_id"`
		Position    vcago.Position `json:"position" bson:"position"`
	}
	CrewList  []Crew
	CrewQuery struct {
		ID    []string `query:"id,omitempty" qs:"id"`
		Name  string   `query:"name" qs:"name"`
		Email string   `query:"email" qs:"email"`
	}
	CrewSimple struct {
		ID    string `json:"id" bson:"id"`
		Name  string `json:"name" bson:"name"`
		Email string `json:"email" bson:"email"`
	}
	CrewParam struct {
		ID string `param:"id"`
	}
)

func (i *CrewCreate) Crew() *Crew {
	return &Crew{
		ID:       uuid.NewString(),
		Name:     i.Name,
		Email:    i.Email,
		Cities:   i.Cities,
		Modified: vcago.NewModified(),
	}
}

func (i *CrewParam) Pipeline() mongo.Pipeline {
	match := vmdb.NewMatch()
	match.EqualString("_id", i.ID)
	return vmdb.NewPipeline().Pipe
}

func (i *CrewQuery) Pipeline() mongo.Pipeline {
	match := vmdb.NewMatch()
	match.EqualStringList("_id", i.ID)
	match.LikeString("email", i.Email)
	match.LikeString("name", i.Name)
	return vmdb.NewPipeline().Match(match).Pipe
}

func CrewPermission(token *vcapool.AccessToken) (err error) {
	if !token.Roles.Validate("employee;admin") {
		return vcago.NewPermissionDenied("crew", nil)
	}
	return
}

func (i *CrewUpdate) Filter() bson.D {
	return bson.D{{Key: "_id", Value: i.ID}}
}

func (i *CrewParam) Filter() bson.D {
	return bson.D{{Key: "_id", Value: i.ID}}
}
