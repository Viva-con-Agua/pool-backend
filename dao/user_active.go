package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

type UserActive vcapool.UserActive

type UserActiveRequest struct {
	vcapool.UserActiveRequest
}

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

//Confirm confirmes UserActive state
func (i *UserActiveRequest) Confirm(ctx context.Context, userID string) (r *vcapool.UserActive, err error) {
	if err = UserActiveCollection.UpdateOneSet(ctx, bson.M{"user_id": i.UserID}, i.Confirmed()); err != nil {
		return
	}
	r = new(vcapool.UserActive)
	err = UserActiveCollection.FindOne(ctx, bson.M{"user_id": i.UserID}, r)
	return
}

//Reject rejects UserActive state
func (i *UserActiveRequest) Reject(ctx context.Context, id string) (r *vcapool.UserActive, err error) {
	if err = UserActiveCollection.UpdateOneSet(ctx, bson.M{"user_id": i.UserID}, i.Rejected()); err != nil {
		return
	}
	r = new(vcapool.UserActive)
	err = UserActiveCollection.FindOne(ctx, bson.M{"user_id": i.UserID}, r)
	return
}

func (i *UserActive) Permission(ctx context.Context, filter bson.M) (err error) {
	err = UserActiveCollection.Permission(ctx, filter, i)
	return
}
