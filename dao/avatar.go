package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

type Avatar vcapool.Avatar

var AvatarCollection = Database.Collection("avatar").CreateIndex("user_id", true)

func (i *Avatar) Create(ctx context.Context) (err error) {
	i.Modified = vcago.NewModified()
	err = AvatarCollection.InsertOne(ctx, i)
	return
}

func (i *Avatar) Delete(ctx context.Context) (err error) {
	i.Modified.Update()
	err = AvatarCollection.DeleteOne(ctx, bson.M{"_id": i.ID})
	return
}
