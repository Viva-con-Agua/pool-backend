package models

import (
	"mime/multipart"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	AvatarCreate struct {
		FileID string `bson:"file_id" json:"file_id"`
		URL    string `bson:"url" json:"url"`
		Type   string `bson:"type" json:"type"`
	}
	Avatar struct {
		ID       string        `bson:"_id" json:"id"`
		FileID   string        `bson:"file_id" json:"file_id"`
		UserID   string        `bson:"user_id" json:"user_id"`
		Modified vmod.Modified `bson:"modified" json:"modified"`
	}
	AvatarParam struct {
		ID string `param:"id"`
	}
	AvatarFile struct {
		File   multipart.File
		Header *multipart.FileHeader
	}
)

var AvatarCollection = "avatar"

func (i *AvatarCreate) Avatar(userID string) *Avatar {
	return &Avatar{
		ID:       uuid.NewString(),
		FileID:   i.FileID,
		UserID:   userID,
		Modified: vmod.NewModified(),
	}
}
func NewAvatar(token *vcapool.AccessToken) *Avatar {
	id := uuid.NewString()
	return &Avatar{
		ID:       id,
		UserID:   token.ID,
		FileID:   id,
		Modified: vmod.NewModified(),
	}
}

type AvatarUpdate struct {
	ID     string `bson:"_id" json:"id"`
	FileID string `bson:"file_id" json:"file_id"`
	URL    string `bson:"url" json:"url"`
	Type   string `bson:"type" json:"type"`
}

func (i *AvatarUpdate) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	filter.EqualString("user_id", token.ID)
	return filter.Bson()
}

func (i *AvatarParam) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	filter.EqualString("user_id", token.ID)
	return filter.Bson()
}

func (i *AvatarParam) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *AvatarParam) FilterChunk() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("files_id", i.ID)
	return filter.Bson()
}
