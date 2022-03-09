package dao

import (
	"context"
	"errors"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

type UserActive vcapool.UserActive

var UserActiveCollection = Database.Collection("user_active").CreateIndex("user_id", true)

func (i *UserActive) Create(ctx context.Context, user *vcapool.User) (r *UserActive, err error) {
	ua := vcapool.NewUserActive(user.ID)
	r = (*UserActive)(ua)
	err = UserActiveCollection.InsertOne(ctx, r)
	return
}

func (i *UserActive) Get(ctx context.Context, filter bson.M) (err error) {
	err = UserActiveCollection.FindOne(ctx, filter, i)
	return
}

func (i *UserActive) Request(ctx context.Context, user *vcapool.User) (err error) {
	if user.Crew.CrewID == "" {
		return vcago.NewStatusBadRequest(errors.New("not an crew member"))
	}
	ua := (*vcapool.UserActive)(i)
	ua.Requested()
	i = (*UserActive)(ua)
	err = UserActiveCollection.UpdateOne(ctx, bson.M{"_id": i.ID}, i)
	return

}

func (i *UserActive) Withdraw(ctx context.Context) (r *UserActive, err error) {
	ua := (*vcapool.UserActive)(i)
	ua.Withdraw()
	r = (*UserActive)(ua)
	err = UserActiveCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, r)
	return
}

type UserActiveRequest struct {
	UserID string `json:"user_id"`
	State  bool   `json:"state"`
}

func (i *UserActiveRequest) Confirm(ctx context.Context) (r *UserActive, err error) {
	userActive := new(vcapool.UserActive)
	if err = UserActiveCollection.FindOne(ctx, bson.M{"user_id": i.UserID}, userActive); err != nil {
		return
	}
	if !userActive.IsRequested() {
		return r, vcago.NewStatusBadRequest(errors.New("active state is not requested"))
	}
	userActive.Confirmed()
	r = (*UserActive)(userActive)
	err = UserActiveCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, r)
	return
}

func (i *UserActive) Reject(ctx context.Context, id string) (r *UserActive, err error) {
	userActive := new(vcapool.UserActive)
	if err = UserActiveCollection.FindOne(ctx, bson.M{"user_id": id}, userActive); err != nil {
		return
	}
	userActive.Rejected()
	r = (*UserActive)(userActive)
	err = UserActiveCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, r)
	return
}

func (i *UserActive) Permission(ctx context.Context, filter bson.M) (err error) {
	err = UserActiveCollection.Permission(ctx, filter, i)
	return
}
