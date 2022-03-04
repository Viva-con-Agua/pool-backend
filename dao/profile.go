package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Profile vcapool.Profile

var ProfilesCollection = Database.Collection("profiles").CreateIndex("user_id", true)

func (i *Profile) Create(ctx context.Context) (err error) {
	i.ID = uuid.NewString()
	i.Modified = vcago.NewModified()
	err = ProfilesCollection.InsertOne(ctx, &i)
	return
}

func (i *Profile) Get(ctx context.Context, id string) (err error) {
	err = ProfilesCollection.FindOne(ctx, bson.M{"_id": id}, &id)
	return
}

func (i *Profile) Update(ctx context.Context) (err error) {
	i.Modified.Update()
	err = ProfilesCollection.UpdateOne(ctx, bson.M{"_id": i.ID}, &i)
	return
}

func (i *Profile) Delete(ctx context.Context) (err error) {
	err = ProfilesCollection.DeleteOne(ctx, bson.M{"_id": i.ID})
	return
}

type ProfileList []Profile

func (i *ProfileList) Get(ctx context.Context, filter bson.M) (err error) {
	err = ProfilesCollection.Find(ctx, filter, &i)
	return
}
