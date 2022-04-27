package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

type AvatarCreate struct {
	vcapool.AvatarCreate
}

type AvatarUpdate struct {
	vcapool.AvatarUpdate
}

type AvatarParam struct {
	vcapool.AvatarParam
}

type Avatar vcapool.Avatar

var AvatarCollection = Database.Collection("avatar").CreateIndex("user_id", true)

func (i *AvatarCreate) Create(ctx context.Context, token *vcapool.AccessToken) (r *vcapool.Avatar, err error) {
	r = i.Avatar(token.ID)
	err = AvatarCollection.InsertOne(ctx, r)
	return
}

func (i *AvatarUpdate) Update(ctx context.Context, token *vcapool.AccessToken) (r *vcapool.Avatar, err error) {
	if err = AvatarCollection.UpdateOneSet(ctx, bson.M{"_id": i.ID, "user_id": token.ID}, i); err != nil {
		return
	}
	r = new(vcapool.Avatar)
	if err = AvatarCollection.FindOne(ctx, bson.M{"_id": i.ID}, r); err != nil {
		return
	}
	return
}

func (i *AvatarParam) Delete(ctx context.Context, token *vcapool.AccessToken) (err error) {
	err = AvatarCollection.DeleteOne(ctx, bson.M{"_id": i.ID, "user_id": token.ID})
	return
}
