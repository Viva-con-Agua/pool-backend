package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

type UserActive vcapool.UserActive

var UserActiveCollection = Database.Collection("user_active").CreateIndex("user_id", true)

func (i *UserActive) Create(ctx context.Context, id string) (r *UserActive, err error) {
	ua := vcapool.NewUserActive(id)
	r = (*UserActive)(ua)
	err = UserActiveCollection.InsertOne(ctx, r)
	return
}

func (i *UserActive) Get(ctx context.Context, filter bson.M) (err error) {
	err = UserActiveCollection.FindOne(ctx, filter, i)
	return
}

func (i *UserActive) Request(ctx context.Context) (err error) {
	ua := (*vcapool.UserActive)(i)
	ua.Requested()
	i = (*UserActive)(ua)
	update := bson.M{"$set": i}
	err = UserActiveCollection.UpdateOne(ctx, bson.M{"_id": i.ID}, update)
	return

}

func (i *UserActive) Withdraw(ctx context.Context) (err error) {
	ua := (*vcapool.UserActive)(i)
	ua.Withdraw()
	i = (*UserActive)(ua)
	update := bson.M{"$set": i}
	err = UserActiveCollection.UpdateOne(ctx, bson.M{"_id": i.ID}, update)
	return
}

type UserActiveRequest struct {
	UserID string `json:"user_id"`
	State  bool   `json:"state"`
}

func (i *UserActive) Confirm(ctx context.Context, userID string) (err error) {
	userActive := new(vcapool.UserActive)
	if err = UserActiveCollection.FindOne(ctx, bson.M{"user_id": userID}, userActive); err != nil {
		return
	}
	if !userActive.IsRequested() {
		return vcago.NewBadRequest("user_active", "active state is not requested")
	}
	userActive.Confirmed()
	i = (*UserActive)(userActive)
	update := bson.M{"$set": i}
	err = UserActiveCollection.UpdateOne(ctx, bson.M{"_id": i.ID}, update)
	return
}

func (i *UserActive) Reject(ctx context.Context, id string) (err error) {
	userActive := new(vcapool.UserActive)
	if err = UserActiveCollection.FindOne(ctx, bson.M{"user_id": id}, userActive); err != nil {
		return
	}
	userActive.Rejected()
	i = (*UserActive)(userActive)
	update := bson.M{"$set": i}
	err = UserActiveCollection.UpdateOne(ctx, bson.M{"_id": i.ID}, update)
	return
}

func (i *UserActive) Permission(ctx context.Context, filter bson.M) (err error) {
	err = UserActiveCollection.Permission(ctx, filter, i)
	return
}
