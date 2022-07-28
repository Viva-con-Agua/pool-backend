package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

type UserInsert struct {
	ID            string         `json:"id,omitempty" bson:"_id"`
	Email         string         `json:"email" bson:"email" validate:"required,email"`
	FirstName     string         `bson:"first_name" json:"first_name" validate:"required"`
	LastName      string         `bson:"last_name" json:"last_name" validate:"required"`
	FullName      string         `bson:"full_name" json:"full_name"`
	DisplayName   string         `bson:"display_name" json:"display_name"`
	Roles         vcago.RoleList `json:"system_roles" bson:"system_roles"`
	Country       string         `bson:"country" json:"country"`
	PrivacyPolicy bool           `bson:"privacy_policy" json:"privacy_policy"`
	Confirmd      bool           `bson:"confirmed" json:"confirmed"`
	LastUpdate    string         `bson:"last_update" json:"last_update"`
	Modified      vcago.Modified `json:"modified" bson:"modified"`
}

type User vcapool.User

type UserDatabase struct {
	vcapool.UserDatabase
}

func GetSendMail(ctx context.Context, currentUser string, contactUser string, scope string) (r *vcago.MailData, err error) {
	user := new(User)
	if err = UserCollection.FindOne(ctx, bson.M{"_id": currentUser}, user); err != nil {
		return
	}
	cUser := new(User)
	if err = UserCollection.FindOne(ctx, bson.M{"_id": contactUser}, cUser); err != nil {
		return
	}
	r = vcago.NewMailData(cUser.Email, "pool-user", scope, cUser.Country)
	r.AddCurrentUser(user.ID, user.Email, user.FirstName, user.LastName)
	r.AddContactUser(cUser.ID, cUser.Email, cUser.FirstName, cUser.LastName)
	return
}

//UseUserCollection represents the user collection
var UserCollection = Database.Collection("users").CreateIndex("email", true)

//Create handles vcago.User model that is providing by auth-service.
func (i *UserDatabase) Create(ctx context.Context) (r *vcapool.User, err error) {
	if err = UserCollection.InsertOne(ctx, &i); err != nil {
		return
	}
	r = i.User()
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
	return
}

//Get selects an User from database
func (i *User) Get(ctx context.Context, filter bson.M) (err error) {
	if err = UserCollection.FindOne(ctx, filter, &i); err != nil {
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
	i.Profile = vcapool.Profile(*profile)
	i.Address = vcapool.Address(*address)
	i.Crew = vcapool.UserCrew(*userCrew)
	i.Active = vcapool.UserActive(*userActive)
	i.NVM = vcapool.UserNVM(*userNVM)
	i.PoolRoles = *poolRoles
	i.Avatar = *avatar

	return
}

func UserDelete(ctx context.Context, token *vcapool.AccessToken) (err error) {
	if err = UserCollection.DeleteOne(ctx, bson.M{"_id": token.ID}); err != nil {
		return
	}
	filter := bson.M{"user_id": token.ID}
	if err = ProfilesCollection.DeleteOne(ctx, filter); err != nil && !vcago.MongoNoDeleted(err) {
		return
	}
	err = nil
	if err = AddressesCollection.DeleteOne(ctx, filter); err != nil && !vcago.MongoNoDeleted(err) {
		return
	}
	err = nil
	if err = UserCrewCollection.DeleteOne(ctx, filter); err != nil && !vcago.MongoNoDeleted(err) {
		return
	}
	err = nil
	if err = UserActiveCollection.DeleteOne(ctx, filter); err != nil && !vcago.MongoNoDeleted(err) {
		return
	}
	err = nil
	if err = UserNVMCollection.DeleteOne(ctx, filter); err != nil && !vcago.MongoNoDeleted(err) {
		return
	}
	err = nil
	if err = PoolRoleCollection.DeleteMany(ctx, filter); err != nil && !vcago.MongoNoDeleted(err) {
		return
	}
	err = nil
	if err = AvatarCollection.DeleteOne(ctx, filter); err != nil && !vcago.MongoNoDeleted(err) {
		return
	}
	err = nil
	return
}

type UserList []vcapool.User

type UserQuery vcapool.UserQuery

func (i *UserQuery) Match() *vcago.MongoMatch {
	match := vcago.NewMongoMatch()
	match.LikeString("first_name", i.FirstName)
	match.LikeString("last_name", i.LastName)
	match.LikeString("full_name", i.FullName)
	match.LikeString("display_name", i.DisplayName)
	match.EqualString("crew.crew_id", i.CrewID)
	match.ElemMatchList("system_roles", "name", i.SystemRoles)
	match.ElemMatchList("pool_roles", "name", i.PoolRoles)
	match.EqualBool("privacy_policy", i.PrivacyPolicy)
	match.StringList("active.status", i.ActiveState)
	match.StringList("nvm.status", i.NVMState)
	match.EqualString("crew.crew_id", i.CrewID)
	match.EqualString("country", i.Country)
	match.EqualBool("confirmed", i.Confirmed)
	match.GteInt64("modified.updated", i.UpdatedFrom)
	match.GteInt64("modified.created", i.CreatedFrom)
	match.LteInt64("modified.updated", i.UpdatedTo)
	match.LteInt64("modified.created", i.CreatedTo)
	return match
}

func (i *UserQuery) List(ctx context.Context) (r *UserList, err error) {
	pipe := vcago.NewMongoPipe()
	pipe.LookupUnwind(AddressesCollection.Name, "_id", "user_id", "address")
	pipe.LookupUnwind(ProfilesCollection.Name, "_id", "user_id", "profile")
	pipe.LookupUnwind(UserCrewCollection.Name, "_id", "user_id", "crew")
	pipe.LookupUnwind(UserActiveCollection.Name, "_id", "user_id", "active")
	pipe.LookupUnwind(UserNVMCollection.Name, "_id", "user_id", "nvm")
	pipe.Lookup(PoolRoleCollection.Name, "_id", "user_id", "pool_roles")
	pipe.LookupUnwind(AvatarCollection.Name, "_id", "user_id", "avatar")
	pipe.Match(i.Match())
	r = new(UserList)
	err = UserCollection.Aggregate(ctx, pipe.Pipe, r)
	return
}
func (i *User) Permission(ctx context.Context, filter bson.M) (err error) {
	err = UserCollection.Permission(ctx, filter, i)
	return
}