package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

/*
type (
	Profile struct {
		ID          string `bson:"_id" json:"id"`
		FirstName   string `bson:"first_name" json:"first_name" validate:"required"`
		LastName    string `bson:"last_name" json:"last_name" validate:"required"`
		FullName    string `bson:"full_name" json:"full_name"`
		DisplayName string `bson:"display_name" json:"display_name"`
		Gender      string `bson:"gender" json:"gender"`
		Avatar      Avatar `bson:"avatar" json:"avatar"`
		UserID      string `bson:"user_id" json:"user_id"`
	}
	Avatar struct {
		URL  string `bson:"url" json:"url"`
		Type string `bson:"type" json:"type"`
	}
)*/

type Profile vcapool.Profile

func NewProfile(user *vcago.User) *Profile {
	return &Profile{
		FirstName: user.Profile.FirstName,
		LastName:  user.Profile.LastName,
		FullName:  user.Profile.FullName,
		UserID:    user.ID,
	}
}

var ProfilesCollection = Database.Collection("profiles")

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
