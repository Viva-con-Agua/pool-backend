package models

import (
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
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
)

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
