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
	AddressCreate struct {
		Street      string `json:"street" bson:"street" validate:"required"`
		Number      string `json:"number" bson:"number" validate:"required"`
		Zip         string `json:"zip" bson:"zip" validate:"required"`
		City        string `json:"city" bson:"city" validate:"required"`
		Country     string `json:"country" bson:"country" validate:"required"`
		CountryCode string `json:"country_code" bson:"country_code" validate:"required"`
		Additional  string `json:"additional" bson:"additional"`
	}
	UsersAddressCreate struct {
		UserID      string `json:"user_id" bson:"user_id" validate:"required"`
		Street      string `json:"street" bson:"street" validate:"required"`
		Number      string `json:"number" bson:"number" validate:"required"`
		Zip         string `json:"zip" bson:"zip" validate:"required"`
		City        string `json:"city" bson:"city" validate:"required"`
		Country     string `json:"country" bson:"country" validate:"required"`
		CountryCode string `json:"country_code" bson:"country_code" validate:"required"`
		Additional  string `json:"additional" bson:"additional"`
	}
	AddressUpdate struct {
		ID          string `json:"id" bson:"_id"`
		Street      string `json:"street" bson:"street"`
		Number      string `json:"number" bson:"number"`
		Zip         string `json:"zip" bson:"zip"`
		City        string `json:"city" bson:"city"`
		Country     string `json:"country" bson:"country"`
		CountryCode string `json:"country_code" bson:"country_code"`
		Additional  string `json:"additional" bson:"additional"`
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
		Additional  string        `json:"additional" bson:"additional"`
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
	AddressImport struct {
		Street      string `json:"street" bson:"street"`
		Number      string `json:"number" bson:"number"`
		Zip         string `json:"zip" bson:"zip"`
		City        string `json:"city" bson:"city"`
		Country     string `json:"country" bson:"country"`
		CountryCode string `json:"country_code" bson:"country_code"`
		Additional  string `json:"additional" bson:"additional"`
		DropsID     string `json:"drops_id"`
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
		CountryCode: i.CountryCode,
		Additional:  i.Additional,
		UserID:      userID,
		Modified:    vmod.NewModified(),
	}
}

func (i *UsersAddressCreate) Address(userID string) (r *Address) {
	return &Address{
		ID:          uuid.NewString(),
		Street:      i.Street,
		Number:      i.Number,
		Zip:         i.Zip,
		City:        i.City,
		Country:     i.Country,
		CountryCode: i.CountryCode,
		Additional:  i.Additional,
		UserID:      userID,
		Modified:    vmod.NewModified(),
	}
}

func (i *AddressImport) Address(userID string) (r *Address) {
	return &Address{
		ID:          uuid.NewString(),
		Street:      i.Street,
		Number:      i.Number,
		Zip:         i.Zip,
		City:        i.City,
		Country:     i.Country,
		CountryCode: i.Country,
		Additional:  i.Additional,
		UserID:      userID,
	}
}

func AddressPermission(token *vcapool.AccessToken) (err error) {
	if !token.Roles.Validate("admin") {
		return vcago.NewPermissionDenied(CrewCollection)
	}
	return
}

func (i *AddressQuery) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	if token.Roles.Validate("employee;admin") {
		filter.EqualStringList("_id", i.ID)
		filter.EqualStringList("crew_id", i.CrewID)
		filter.EqualStringList("user_id", i.UserID)
	} else {
		filter.EqualString("user_id", token.ID)
	}
	filter.GteInt64("modified.updated", i.UpdatedFrom)
	filter.GteInt64("modified.created", i.CreatedFrom)
	filter.LteInt64("modified.updated", i.UpdatedTo)
	filter.LteInt64("modified.created", i.CreatedTo)
	return filter.Bson()
}

func (i *AddressUpdate) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	filter.EqualString("user_id", token.ID)
	return filter.Bson()
}

func (i *AddressUpdate) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *AddressParam) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	filter.EqualString("user_id", token.ID)
	return filter.Bson()
}

func (i *AddressParam) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *AddressImport) FilterUser() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("drops_id", i.DropsID)
	return filter.Bson()
}
