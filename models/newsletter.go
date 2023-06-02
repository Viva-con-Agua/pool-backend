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
	NewsletterCreate struct {
		Value  string `json:"value" bson:"value"`
		UserID string `json:"user_id" bson:"user_id"`
	}
	Newsletter struct {
		ID       string        `json:"id" bson:"_id"`
		Value    string        `json:"value" bson:"value"`
		UserID   string        `json:"user_id" bson:"user_id"`
		Modified vmod.Modified `json:"modified" bson:"modified"`
	}
	NewsletterParam struct {
		ID string `param:"id"`
	}
	NewsletterImport struct {
		Value   string `json:"value"`
		DropsID string `json:"drops_id"`
	}
)

var NewsletterCollection = "newsletters"

func NewsletterDeletePermission(token *vcapool.AccessToken) (err error) {
	if !token.Roles.Validate("employee;admin") {
		return vcago.NewPermissionDenied(ArtistCollection)
	}
	return
}

func (i *NewsletterCreate) Newsletter(token *vcapool.AccessToken) *Newsletter {
	return &Newsletter{
		ID:       uuid.NewString(),
		Value:    i.Value,
		UserID:   token.ID,
		Modified: vmod.NewModified(),
	}
}

func (i *NewsletterCreate) NewsletterAdmin() *Newsletter {
	return &Newsletter{
		ID:       uuid.NewString(),
		Value:    i.Value,
		UserID:   i.UserID,
		Modified: vmod.NewModified(),
	}
}

func (i *NewsletterImport) ToNewsletter(userID string) *Newsletter {
	return &Newsletter{
		ID:       uuid.NewString(),
		Value:    i.Value,
		UserID:   userID,
		Modified: vmod.NewModified(),
	}
}

func (i *NewsletterParam) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *Newsletter) DeletePermission(token *vcapool.AccessToken) (err error) {
	if !token.Roles.Validate("employee;admin") && token.ID != i.UserID {
		return vcago.NewPermissionDenied(NewsletterCollection)
	}
	return
}
