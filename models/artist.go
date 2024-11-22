package models

import (
	"strings"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	ArtistCreate struct {
		Name string `json:"name" bson:"name" validate:"required"`
	}
	Artist struct {
		ID       string        `json:"id" bson:"_id"`
		Name     string        `json:"name" bson:"name"`
		Modified vmod.Modified `json:"modified" bson:"modified"`
	}
	ArtistUpdate struct {
		ID   string `json:"id" bson:"_id"`
		Name string `json:"name" bson:"name"`
	}
	ArtistParam struct {
		ID string `param:"id"`
	}
	ArtistQuery struct {
		ID          string `query:"id" qs:"id"`
		Name        string `query:"name" qs:"name"`
		UpdatedTo   string `query:"updated_to" qs:"updated_to"`
		UpdatedFrom string `query:"updated_from" qs:"updated_from"`
		CreatedTo   string `query:"created_to" qs:"created_to"`
		CreatedFrom string `query:"created_from" qs:"created_from"`
	}
)

var ArtistCollection = "artists"

func ArtistPermission(token *AccessToken) (err error) {
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPEventRole)) {
		return vcago.NewPermissionDenied(ArtistCollection)
	}
	return
}

func ArtistDeletePermission(token *AccessToken) (err error) {
	if !token.Roles.Validate("admin;employee;pool_employee") {
		return vcago.NewPermissionDenied(ArtistCollection)
	}
	return
}

func ToArtistList(artists []Artist) string {
	names := make([]string, len(artists))
	for i, artist := range artists {
		names[i] = artist.Name
	}

	result := strings.Join(names, ", ")
	return result
}

func (i *ArtistCreate) Artist() *Artist {
	return &Artist{
		ID:       uuid.NewString(),
		Name:     i.Name,
		Modified: vmod.NewModified(),
	}
}

func (i *ArtistParam) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *ArtistUpdate) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *ArtistQuery) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	filter.LikeString("name", i.Name)
	filter.GteInt64("modified.updated", i.UpdatedFrom)
	filter.GteInt64("modified.created", i.CreatedFrom)
	filter.LteInt64("modified.updated", i.UpdatedTo)
	filter.LteInt64("modified.created", i.CreatedTo)
	return filter.Bson()
}
