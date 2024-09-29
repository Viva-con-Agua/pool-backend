package models

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	OrganisationCreate struct {
		Name         string `json:"name" bson:"name" validate:"required"`
		Abbreviation string `json:"abbreviation" bson:"abbreviation"`
		Email        string `json:"email" bson:"email"`
	}
	Organisation struct {
		ID           string        `json:"id" bson:"_id"`
		Name         string        `json:"name" bson:"name"`
		Abbreviation string        `json:"abbreviation" bson:"abbreviation"`
		Email        string        `json:"email" bson:"email"`
		Modified     vmod.Modified `json:"modified" bson:"modified"`
	}
	OrganisationUpdate struct {
		ID           string `json:"id" bson:"_id"`
		Abbreviation string `json:"abbreviation" bson:"abbreviation"`
		Email        string `json:"email" bson:"email"`
		Name         string `json:"name" bson:"name"`
	}
	OrganisationParam struct {
		ID string `param:"id"`
	}
	OrganisationQuery struct {
		ID           string `query:"id" qs:"id"`
		Name         string `query:"name" qs:"name"`
		Abbreviation string `query:"abbreviation" qs:"abbreviation"`
		Email        string `query:"email" qs:"email"`
		UpdatedTo    string `query:"updated_to" qs:"updated_to"`
		UpdatedFrom  string `query:"updated_from" qs:"updated_from"`
		CreatedTo    string `query:"created_to" qs:"created_to"`
		CreatedFrom  string `query:"created_from" qs:"created_from"`
	}
)

var OrganisationCollection = "organisations"

func OrganisationPermission(token *vcapool.AccessToken) (err error) {
	if !token.Roles.Validate("admin") {
		return vcago.NewPermissionDenied(OrganisationCollection)
	}
	return
}

func (i *OrganisationCreate) Organisation() *Organisation {
	return &Organisation{
		ID:           uuid.NewString(),
		Name:         i.Name,
		Email:        i.Email,
		Abbreviation: i.Abbreviation,
		Modified:     vmod.NewModified(),
	}
}

func (i *OrganisationParam) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *OrganisationUpdate) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *OrganisationQuery) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	filter.LikeString("name", i.Name)
	filter.GteInt64("modified.updated", i.UpdatedFrom)
	filter.GteInt64("modified.created", i.CreatedFrom)
	filter.LteInt64("modified.updated", i.UpdatedTo)
	filter.LteInt64("modified.created", i.CreatedTo)
	return filter.Bson()
}
