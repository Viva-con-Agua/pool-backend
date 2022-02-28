package dao

import (
	"context"
	"errors"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
)

type UserActive vcapool.UserActive

var UserActiveCollection = Database.Collection("user_active").CreateIndex("user_id", true)

func (i *UserActive) Create(ctx context.Context, user *vcapool.User) (r *UserActive, err error) {
	if user.Crew.CrewID == "" {
		return nil, vcago.NewStatusBadRequest(errors.New("not an crew member"))
	}
	ua := vcapool.NewUserActive(user.ID)
	r = (*UserActive)(ua)
	err = UserActiveCollection.InsertOne(ctx, r)
	return
}
