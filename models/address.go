package models

import (
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	AddressCreate struct {
		Street      string `json:"street" bson:"street" validate:"required"`
		Number      string `json:"number" bson:"number" validate:"required"`
		Zip         string `json:"zip" bson:"zip" validate:"required"`
		City        string `json:"city" bson:"city" validate:"required"`
		Country     string `json:"country" bson:"country" validate:"required"`
		CountryCode string `json:"country_code" bson:"country_code" validate:"required"`
		Additionals string `json:"additionals" bson:"additionals"`
	}
	AddressUpdate struct {
		ID          string `json:"id" bson:"_id"`
		Street      string `json:"street" bson:"street"`
		Number      string `json:"number" bson:"number"`
		Zip         string `json:"zip" bson:"zip"`
		City        string `json:"city" bson:"city"`
		Country     string `json:"country" bson:"country"`
		CountryCode string `json:"country_code" bson:"country_code"`
		Additionals string `json:"additionals" bson:"additionals"`
		UserID      string `json:"user_id" bson:"user_id"`
	}
	Address struct {
		ID          string        `json:"id" bson:"_id"`
		Street      string        `json:"street" bson:"street"`
		Number      string        `json:"number" bson:"number"`
		Zip         string        `json:"zip" bson:"zip"`
		City        string        `json:"city" bson:"city"`
		Country     string        `json:"country" bson:"country"`
		CountryCode string        `json:"country_code" bson:"country_code"`
		Additionals string        `json:"additionals" bson:"additionals"`
		UserID      string        `json:"user_id" bson:"user_id"`
		Modified    vmod.Modified `json:"modified" bson:"modified"`
	}
	AddressQuery struct {
		ID          []string `query:"id" qs:"id"`
		CrewID      []string `query:"crew_id" qs:"crew_id"`
		UserID      []string `query:"user_id" qs:"user_id"`
		UpdatedTo   string   `query:"updated_to" qs:"updated_to"`
		UpdatedFrom string   `query:"updated_from" qs:"updated_from"`
		CreatedTo   string   `query:"created_to" qs:"created_to"`
		CreatedFrom string   `query:"created_from" qs:"created_from"`
	}
	AddressParam struct {
		ID string `param:"id"`
	}
)

var AddressesCollection = "addresses"

func (i *AddressCreate) Address(userID string) (r *Address) {
	return &Address{
		ID:          uuid.NewString(),
		Street:      i.Street,
		Number:      i.Number,
		Zip:         i.Zip,
		City:        i.City,
		Country:     i.Country,
		CountryCode: i.Country,
		Additionals: i.Additionals,
		UserID:      userID,
	}
}

func (i *AddressParam) Pipeline(token *vcapool.AccessToken) mongo.Pipeline {
	match := vmdb.NewFilter()
	match.EqualString("_id", i.ID)
	if !token.Roles.Validate("employee;admin") {
		match.EqualString("user_id", token.ID)
	}
	return vmdb.NewPipeline().Match(match.Bson()).Pipe
}

func (i *AddressQuery) Filter(token *vcapool.AccessToken) bson.D {
	match := vmdb.NewFilter()
	if token.Roles.Validate("employee;admin") {
		match.EqualStringList("_id", i.ID)
		match.EqualStringList("crew_id", i.CrewID)
		match.EqualStringList("user_id", i.UserID)
	} else {
		match.EqualString("user_id", token.ID)
	}
	match.GteInt64("modified.updated", i.UpdatedFrom)
	match.GteInt64("modified.created", i.CreatedFrom)
	match.LteInt64("modified.updated", i.UpdatedTo)
	match.LteInt64("modified.created", i.CreatedTo)
	return match.Bson()
}

func (i *AddressUpdate) Filter(token *vcapool.AccessToken) bson.D {
	return bson.D{{Key: "_id", Value: i.ID}, {Key: "user_id", Value: token.ID}}
}

func (i *AddressParam) Filter(token *vcapool.AccessToken) bson.D {
	return bson.D{{Key: "_id", Value: i.ID}, {Key: "user_id", Value: token.ID}}
}
