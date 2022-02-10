package dao

import (
	"context"
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type User vcapool.User

//NewUser creates an new User from given vcago.User
func NewUser(user *vcago.User) (r *User) {
	return &User{
		ID:       user.ID,
		Email:    user.Email,
		Profile:  user.Profile,
		Roles:    user.Roles,
		Modified: vcago.NewModified(),
	}
}

func (i *User) ToAuthToken() (r *AuthToken, err error) {
	token := new(vcapool.AuthToken)
	if token, err = vcapool.NewAuthToken(i.ToVPUser()); err != nil {
		return
	}
	return &AuthToken{
		ID:        uuid.NewString(),
		Token:     *token,
		UserID:    i.ID,
		ExpiresAt: time.Now().Unix(),
	}, nil
}

//ToVPUser converts an User to vpool.User
func (i *User) ToVPUser() *vcapool.User {
	return &vcapool.User{
		ID:       i.ID,
		Email:    i.Email,
		Profile:  i.Profile,
		CrewID:   i.CrewID,
		Crew:     i.Crew,
		Address:  i.Address,
		Roles:    i.Roles,
		Modified: i.Modified,
	}
}

var UserCollection = Database.Collection("users").CreateIndex("email", true)

//Create handles vcago.User model that is providing by auth-service.
func (i *User) Create(ctx context.Context) (err error) {
	err = UserCollection.InsertOne(ctx, &i)
	return
}

//Get selects an User from database
func (i *User) Get(ctx context.Context, id string) (err error) {
	err = UserCollection.FindOne(ctx, bson.M{"_id": id}, &i)
	return
}
