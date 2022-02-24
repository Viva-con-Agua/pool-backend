package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

type UserInsert struct {
	ID       string         `json:"id,omitempty" bson:"_id"`
	Email    string         `json:"email" bson:"email" validate:"required,email"`
	Roles    vcago.RoleList `json:"roles" bson:"roles"`
	Country  string         `json:"country" bson:"country"`
	Modified vcago.Modified `json:"modified" bson:"modified"`
}

type User vcapool.User

//NewUser creates an new User from given vcago.User
func NewUser(user *vcago.User) (r *UserInsert) {
	return &UserInsert{
		ID:       user.ID,
		Email:    user.Email,
		Roles:    user.Roles,
		Country:  user.Country,
		Modified: vcago.NewModified(),
	}
}

//CreateUserFromToken creates user form vcago.User model
func CreateUserFromToken(ctx context.Context, user *vcago.User) (r *User, err error) {
	userInsert := NewUser(user)
	if err = userInsert.Create(ctx); err != nil {
		return
	}
	profile := NewProfile(user)
	if err = profile.Create(ctx); err != nil {
		return
	}
	return &User{
		ID:       userInsert.ID,
		Email:    userInsert.Email,
		Roles:    userInsert.Roles,
		Country:  userInsert.Country,
		Modified: userInsert.Modified,
		Profile:  vcapool.Profile(*profile),
	}, nil
}

//UseUserCollection represents the user collection
var UserCollection = Database.Collection("users").CreateIndex("email", true)

//Create handles vcago.User model that is providing by auth-service.
func (i *UserInsert) Create(ctx context.Context) (err error) {
	err = UserCollection.InsertOne(ctx, &i)
	return
}

//Get selects an User from database
func (i *User) Get(ctx context.Context, id string) (err error) {
	if err = UserCollection.FindOne(ctx, bson.M{"_id": id}, &i); err != nil {
		return
	}
	profile := new(Profile)
	err = ProfilesCollection.FindOne(ctx, bson.M{"user_id": id}, &profile)
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	address := new(Address)
	err = AddressesCollection.FindOne(ctx, bson.M{"user_id": id}, &address)
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	err = nil
	i.Profile = vcapool.Profile(*profile)
	i.Address = vcapool.Address(*address)
	return
}

type UserList []User

func (i *UserList) List(ctx context.Context) (err error) {
	pipe := vcago.NewMongoPipe()
	pipe.AddModelAt("addresses", "_id", "user_id", "address")
	err = UserCollection.Aggregate(ctx, pipe.Pipe, i)
	return
}
