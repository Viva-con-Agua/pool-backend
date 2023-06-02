package dao

/*
import (
	"context"
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	UserMigrate struct {
		ID        string                `json:"id" bson:"_id"`
		Email     string                `json:"email" bson:"email"`
		FirstName string                `json:"first_name" bson:"first_name"`
		LastName  string                `json:"last_name" bson:"last_name"`
		FullName  string                `json:"full_name" bson:"full_name"`
		DropsID   string                `json:"drops_id" bson:"drops_id"`
		PoolID    string                `json:"pool_id" bson:"pool_id"`
		Crew      CrewMigrate           `json:"crew" bson:"crew"`
		Active    string                `json:"active" bson:"active"`
		NVMDate   int64                 `json:"nvm_state" bson:"nvm_state"`
		Address   vcapool.AddressCreate `json:"address" bson:"address"`
		Profile   vcapool.ProfileCreate `json:"profile" bson:"profile"`
		Roles     []string              `json:"roles" bson:"roles"`
		Modified  vcago.Modified        `json:"modified" bson:"modified"`
	}
	CrewMigrate struct {
		ID       string         `json:"id" bson:"id"`
		Modified vcago.Modified `json:"modified" bson:"modified"`
	}
	UserUpdate struct {
		DropsID  string         `json:"drops_id" bson:"drops_id"`
		Modified vcago.Modified `json:"modified" bson:"modified"`
	}
)

func (i *UserMigrate) UserUpdate() *UserUpdate {
	return &UserUpdate{
		DropsID:  i.DropsID,
		Modified: i.Modified,
	}
}

func (i *UserMigrate) UserActive() *vcapool.UserActiveUpdate {
	if i.Active == "active" {
		return &vcapool.UserActiveUpdate{
			Status: "confirmed",
			Since:  time.Now().Unix(),
		}
	} else if i.Active == "requested" {
		return &vcapool.UserActiveUpdate{
			Status: "requested",
			Since:  time.Now().Unix(),
		}
	}
	return nil
}

func (i *UserMigrate) UserNVM(userID string) *vcapool.UserNVM {
	if i.NVMDate != 0 {
		return &vcapool.UserNVM{
			ID:       uuid.NewString(),
			Status:   "confirmed",
			Since:    i.NVMDate,
			Expired:  0,
			UserID:   userID,
			Modified: vcago.NewModified(),
		}
	}
	return nil
}

func (i *UserMigrate) MigrateUser(ctx context.Context) (err error) {
	user := new(vcapool.UserDatabase)
	if err = UserCollection.FindOne(ctx, bson.M{"email": i.Email}, user); err != nil {
		return
	}
	update := bson.M{"$set": i.UserUpdate()}
	if err = UserCollection.UpdateOne(ctx, bson.M{"_id": user.ID}, update); err != nil && !vcago.MongoNoUpdated(err) {
		return
	}
	crew := new(vcapool.Crew)
	if err = CrewsCollection.FindOne(ctx, bson.M{"_id": i.Crew.ID}, crew); err != nil {
		return
	}
	err = nil
	userCrew := vcapool.NewUserCrew(user.ID, i.Crew.ID, crew.Name, crew.Email)
	uc := (*UserCrew)(userCrew)
	uc.Modified = i.Crew.Modified

	if err = UserCrewCollection.InsertOne(ctx, uc); err != nil && !vcago.MongoConfict(err) {
		return
	}
	err = nil
	address := i.Address.Address(user.ID)
	if err = AddressesCollection.InsertOne(ctx, address); err != nil && !vcago.MongoConfict(err) {
		return
	}
	err = nil
	profile := i.Profile.Profile(user.ID)
	if err = ProfileCollection.InsertOne(ctx, profile); err != nil && !vcago.MongoConfict(err) {
		return
	}
	err = nil
	userActive := i.UserActive()
	if userActive != nil {
		if err = UserActiveCollection.UpdateOneSet(ctx, bson.M{"user_id": user.ID}, userActive); err != nil {
			return
		}
	}
	userNVM := i.UserNVM(user.ID)
	if userNVM != nil {
		if err = UserNVMCollection.UpdateOneSet(ctx, bson.M{"user_id": user.ID}, userNVM); err != nil {
			return
		}
	}
	for n := range i.Roles {
		role := new(vcago.Role)
		role, err = vcapool.NewRole(i.Roles[n], user.ID)
		if err = PoolRoleCollection.InsertOne(ctx, role); err != nil && !vcago.MongoConfict(err) {
			return
		}
		err = nil
	}
	return
}
*/
