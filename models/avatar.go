package models

import (
	"mime/multipart"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	Avatar struct {
		ID       string        `bson:"_id" json:"id"`
		FileID   string        `bson:"file_id" json:"file_id"`
		UserID   string        `bson:"user_id" json:"user_id"`
		Modified vmod.Modified `bson:"modified" json:"modified"`
	}
	AvatarUpdate struct {
		ID     string `bson:"_id" json:"id"`
		FileID string `bson:"file_id" json:"file_id"`
		URL    string `bson:"url" json:"url"`
		Type   string `bson:"type" json:"type"`
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
var FSChunkCollection = "fs.chunks"
var FSFilesCollection = "fs.files"

func NewAvatar(token *AccessToken) *Avatar {
	id := uuid.NewString()
	return &Avatar{
		ID:       id,
		UserID:   token.ID,
		FileID:   id,
		Modified: vmod.NewModified(),
	}
}

func NewAvatarClean() *Avatar {
	return &Avatar{}
}

func (i *AvatarUpdate) PermittedFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", token.ID)
	return filter.Bson()
}

func (i *AvatarParam) PermittedFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	filter.EqualString("user_id", token.ID)
	return filter.Bson()
}

func (i *AvatarParam) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *AvatarParam) MatchChunk() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("files_id", i.ID)
	return filter.Bson()
}
