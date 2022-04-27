package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

//UserSystem represents the dao for vcago.User.
type UserSystem struct {
	vcago.User
}

func NewUserSystem(user *vcago.User) *UserSystem {
	return &UserSystem{
		*user,
	}
}

//Create handles vcago.User model that is providing by auth-service.
func (i *UserSystem) Create(ctx context.Context) (r *vcapool.User, err error) {
	database := vcapool.NewUserDatabase(i.User)
	if err = UserCollection.InsertOne(ctx, &database); err != nil {
		return
	}
	r = database.User()
	nvm := vcapool.NewUserNVM(i.ID)
	if err = UserNVMCollection.InsertOne(ctx, nvm); err != nil {
		return
	}
	active := vcapool.NewUserActive(i.ID)
	if err = UserActiveCollection.InsertOne(ctx, active); err != nil {
		return
	}
	r.NVM = *nvm
	r.Active = *active
	vcago.Nats.Publish("user.created", r)
	return
}

//Get get vcapool.User
func (i *UserSystem) Get(ctx context.Context) (r *vcapool.User, err error) {
	r = new(vcapool.User)
	if err = UserCollection.FindOne(ctx, bson.M{"_id": i.ID}, &r); err != nil {
		return
	}
	profile := new(Profile)
	err = ProfilesCollection.FindOne(ctx, bson.M{"user_id": i.ID}, profile)
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	address := new(Address)
	err = AddressesCollection.FindOne(ctx, bson.M{"user_id": i.ID}, address)
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	userCrew := new(UserCrew)
	err = UserCrewCollection.FindOne(ctx, bson.M{"user_id": i.ID}, userCrew)
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	userActive := new(UserActive)
	err = UserActiveCollection.FindOne(ctx, bson.M{"user_id": i.ID}, userActive)
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	userNVM := new(UserNVM)
	err = UserNVMCollection.FindOne(ctx, bson.M{"user_id": i.ID}, userNVM)
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	poolRoles := new(vcago.RoleList)
	if err = PoolRoleCollection.Find(ctx, bson.M{"user_id": i.ID}, poolRoles); err != nil {
		return
	}
	avatar := new(vcapool.Avatar)
	err = AvatarCollection.FindOne(ctx, bson.M{"user_id": i.ID}, avatar)
	if err != nil && !vcago.MongoNoDocuments(err) {
		return
	}
	err = nil
	r.Profile = vcapool.Profile(*profile)
	r.Address = vcapool.Address(*address)
	r.Crew = vcapool.UserCrew(*userCrew)
	r.Active = vcapool.UserActive(*userActive)
	r.NVM = vcapool.UserNVM(*userNVM)
	r.PoolRoles = *poolRoles
	r.Avatar = *avatar
	return
}

func (i *UserSystem) Update(ctx context.Context) (r *vcapool.User, err error) {
	if err = UserCollection.UpdateOneSet(ctx, bson.M{"_id": i.ID}, i); err != nil {
		return
	}
	if r, err = i.Get(ctx); err != nil {
		return
	}
	vcago.Nats.Publish("user.updated", r)
	return
}
