package models

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	CrewCreate struct {
		Name         string `json:"name" bson:"name"`
		Email        string `json:"email" bson:"email"`
		Abbreviation string `json:"abbreviation" bson:"abbreviation"`
		Mattermost   string `bson:"mattermost_username" json:"mattermost_username"`
		Additional   string `json:"additional" bson:"additional"`
		Cities       []City `json:"cities" bson:"cities"`
		Status       string `json:"status" bson:"status"`
	}
	CrewUpdate struct {
		ID           string `json:"id,omitempty" bson:"_id"`
		Name         string `json:"name" bson:"name"`
		Email        string `json:"email" bson:"email"`
		Abbreviation string `json:"abbreviation" bson:"abbreviation"`
		Mattermost   string `bson:"mattermost_username" json:"mattermost_username"`
		Additional   string `json:"additional" bson:"additional"`
		Status       string `json:"status" bson:"status"`
		Cities       []City `json:"cities" bson:"cities"`
	}
	CrewUpdateASP struct {
		ID         string `json:"id,omitempty" bson:"_id"`
		Mattermost string `bson:"mattermost_username" json:"mattermost_username"`
		Additional string `json:"additional" bson:"additional"`
	}
	Crew struct {
		ID           string        `json:"id,omitempty" bson:"_id"`
		Name         string        `json:"name" bson:"name"`
		Email        string        `json:"email" bson:"email"`
		Abbreviation string        `json:"abbreviation" bson:"abbreviation"`
		Mattermost   string        `bson:"mattermost_username" json:"mattermost_username"`
		Additional   string        `json:"additional" bson:"additional"`
		MailboxID    string        `json:"mailbox_id" bson:"mailbox_id"`
		Cities       []City        `json:"cities" bson:"cities"`
		Status       string        `json:"status" bson:"status"`
		Modified     vmod.Modified `json:"modified" bson:"modified"`
	}
	City struct {
		City        string        `json:"city" bson:"city"`
		Country     string        `json:"country" bson:"country"`
		CountryCode string        `json:"country_code" bson:"country_code"`
		PlaceID     string        `json:"place_id" bson:"place_id"`
		Position    vmod.Position `json:"position" bson:"position"`
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

var CrewCollection = "crews"

func (i *CrewCreate) Crew() *Crew {
	return &Crew{
		ID:           uuid.NewString(),
		Name:         i.Name,
		Email:        i.Email,
		Mattermost:   i.Mattermost,
		Abbreviation: i.Abbreviation,
		Additional:   i.Additional,
		Cities:       i.Cities,
		Status:       i.Status,
		Modified:     vmod.NewModified(),
	}
}

func (i *CrewUpdate) ToCrewUpdateASP() *CrewUpdateASP {
	return &CrewUpdateASP{
		ID:         i.ID,
		Mattermost: i.Mattermost,
		Additional: i.Additional,
	}
}

func (i *CrewParam) Pipeline() mongo.Pipeline {
	match := vmdb.NewFilter()
	match.EqualString("_id", i.ID)
	return vmdb.NewPipeline().Match(match.Bson()).Pipe
}

func (i *CrewQuery) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualStringList("_id", i.ID)
	filter.LikeString("email", i.Email)
	filter.LikeString("name", i.Name)
	return bson.D(*filter)
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
