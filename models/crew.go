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
	CrewCreate struct {
		Name           string `json:"name" bson:"name"`
		Email          string `json:"email" bson:"email"`
		Abbreviation   string `json:"abbreviation" bson:"abbreviation"`
		Mattermost     string `bson:"mattermost_username" json:"mattermost_username"`
		Additional     string `json:"additional" bson:"additional"`
		Cities         []City `json:"cities" bson:"cities"`
		OrganisationID string `json:"organisation_id" bson:"organisation_id"`
		Status         string `json:"status" bson:"status"`
		AspSelection   string `json:"asp_selection" bson:"asp_selection"`
	}
	CrewUpdate struct {
		ID             string `json:"id,omitempty" bson:"_id"`
		Name           string `json:"name" bson:"name"`
		Email          string `json:"email" bson:"email"`
		Abbreviation   string `json:"abbreviation" bson:"abbreviation"`
		Mattermost     string `bson:"mattermost_username" json:"mattermost_username"`
		OrganisationID string `json:"organisation_id" bson:"organisation_id"`
		Additional     string `json:"additional" bson:"additional"`
		Status         string `json:"status" bson:"status"`
		AspSelection   string `json:"asp_selection" bson:"asp_selection"`
		Cities         []City `json:"cities" bson:"cities"`
	}
	CrewUpdateASP struct {
		ID         string `json:"id,omitempty" bson:"_id"`
		Mattermost string `bson:"mattermost_username" json:"mattermost_username"`
		Additional string `json:"additional" bson:"additional"`
	}
	Crew struct {
		ID             string        `json:"id,omitempty" bson:"_id"`
		Name           string        `json:"name" bson:"name"`
		Email          string        `json:"email" bson:"email"`
		Abbreviation   string        `json:"abbreviation" bson:"abbreviation"`
		Mattermost     string        `bson:"mattermost_username" json:"mattermost_username"`
		Additional     string        `json:"additional" bson:"additional"`
		OrganisationID string        `json:"organisation_id" bson:"organisation_id"`
		Organisation   Organisation  `json:"organisation" bson:"organisation"`
		MailboxID      string        `json:"mailbox_id" bson:"mailbox_id"`
		Cities         []City        `json:"cities" bson:"cities"`
		Status         string        `json:"status" bson:"status"`
		AspSelection   string        `json:"asp_selection" bson:"asp_selection"`
		Modified       vmod.Modified `json:"modified" bson:"modified"`
	}
	CrewPublic struct {
		ID             string       `json:"id,omitempty" bson:"_id"`
		Name           string       `json:"name" bson:"name"`
		Cities         []City       `json:"cities" bson:"cities"`
		Organisation   Organisation `json:"organisation" bson:"organisation"`
		OrganisationID string       `json:"organisation_id" bson:"organisation_id"`
		Mattermost     string       `bson:"mattermost_username" json:"mattermost_username"`
	}
	CrewName struct {
		ID   string `json:"id,omitempty" bson:"_id"`
		Name string `json:"name" bson:"name"`
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
		ID             []string `query:"id,omitempty" qs:"id"`
		Name           string   `query:"name" qs:"name"`
		Status         string   `json:"status" qs:"status"`
		Organisation   string   `json:"organisation_name" qs:"organisation_name"`
		OrganisationID []string `json:"organisation_id" qs:"organisation_id"`
		Email          string   `query:"email" qs:"email"`
	}
	CrewSimple struct {
		ID           string       `json:"id" bson:"id"`
		Name         string       `json:"name" bson:"name"`
		Email        string       `json:"email" bson:"email"`
		Organisation Organisation `json:"organisation" bson:"organisation"`
	}
	CrewParam struct {
		ID string `param:"id"`
	}
)

var CrewCollection = "crews"

func CrewPermission(token *vcapool.AccessToken) (err error) {
	if !token.Roles.Validate("admin;employee;pool_employee") {
		return vcago.NewPermissionDenied(CrewCollection)
	}
	return
}
func CrewPipeline() *vmdb.Pipeline {
	pipe := vmdb.NewPipeline()
	pipe.LookupUnwind(OrganisationCollection, "organisation_id", "_id", "organisation")
	return pipe
}

func CrewUpdatePermission(token *vcapool.AccessToken) (err error) {
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPRole)) {
		return vcago.NewPermissionDenied(CrewCollection)
	}
	return
}

func (i *CrewCreate) Crew() *Crew {
	return &Crew{
		ID:             uuid.NewString(),
		Name:           i.Name,
		Email:          i.Email,
		Mattermost:     i.Mattermost,
		Abbreviation:   i.Abbreviation,
		Additional:     i.Additional,
		OrganisationID: i.OrganisationID,
		Cities:         i.Cities,
		Status:         i.Status,
		Modified:       vmod.NewModified(),
	}
}

func (i *CrewUpdate) ToCrewUpdateASP() *CrewUpdateASP {
	return &CrewUpdateASP{
		ID:         i.ID,
		Mattermost: i.Mattermost,
		Additional: i.Additional,
	}
}

func (i *CrewQuery) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualStringList("_id", i.ID)
	filter.EqualStringList("organisation_id", i.OrganisationID)
	filter.LikeString("email", i.Email)
	filter.LikeString("status", i.Status)
	filter.LikeString("name", i.Name)
	return filter.Bson()
}

func (i *CrewQuery) ActiveFilter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualStringList("_id", i.ID)
	filter.LikeString("email", i.Email)
	filter.LikeString("status", "active")
	filter.LikeString("name", i.Name)
	return filter.Bson()
}

func (i *CrewQuery) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", token.CrewID)
	filter.LikeString("status", "active")
	return filter.Bson()
}

func (i *CrewUpdate) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	if !token.Roles.Validate("admin;employee;pool_employee") {
		filter.EqualString("_id", token.CrewID)
	} else {
		filter.EqualString("_id", i.ID)
	}
	return filter.Bson()
}

func (i *CrewParam) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	if !token.Roles.Validate("admin;employee;pool_employee") {
		filter.EqualString("_id", token.CrewID)
	} else {
		filter.EqualString("_id", i.ID)
	}
	return filter.Bson()
}

func (i *CrewParam) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}
