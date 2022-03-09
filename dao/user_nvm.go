package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

type UserNVM vcapool.UserNVM

var UserNVMCollection = Database.Collection("user_nvm").CreateIndex("user_id", true)

func (i *UserNVM) Create(ctx context.Context, userID string) (r *UserNVM, err error) {
	ua := vcapool.NewUserNVM(userID)
	r = (*UserNVM)(ua)
	err = UserNVMCollection.InsertOne(ctx, r)
	return
}

func (i *UserNVM) Get(ctx context.Context, filter bson.M) (err error) {
	err = UserNVMCollection.FindOne(ctx, filter, i)
	return
}

func (i *UserNVM) Withdraw(ctx context.Context) (r *UserNVM, err error) {
	ua := (*vcapool.UserNVM)(i)
	ua.Withdraw()
	r = (*UserNVM)(ua)
	err = UserNVMCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, r)
	return
}

func (i *UserNVM) Confirm(ctx context.Context) (r *UserNVM, err error) {
	ua := (*vcapool.UserNVM)(i)
	ua.Confirmed()
	r = (*UserNVM)(ua)
	err = UserNVMCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, r)
	return
}

type UserNVMRequest struct {
	UserID string `json:"user_id"`
	State  bool   `json:"state"`
}

func (i *UserNVM) Reject(ctx context.Context, id string) (r *UserNVM, err error) {
	userNVM := new(vcapool.UserNVM)
	if err = UserNVMCollection.FindOne(ctx, bson.M{"user_id": id}, userNVM); err != nil {
		return
	}
	userNVM.Rejected()
	r = (*UserNVM)(userNVM)
	err = UserNVMCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, r)
	return
}

func (i *UserNVM) Permission(ctx context.Context, filter bson.M) (err error) {
	err = UserNVMCollection.Permission(ctx, filter, i)
	return
}
