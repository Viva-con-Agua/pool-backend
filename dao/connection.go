package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"go.mongodb.org/mongo-driver/bson"
)

var Database *vmdb.Database

// UserCollection represents the database collection of the User model.
var UserCollection *vmdb.Collection

// UserCrewCollection represents the database collection of the UserCrew model.
var UserCrewCollection *vmdb.Collection

// ActiveCollection represents the database collection of the Active model.
var ActiveCollection *vmdb.Collection

// NVMCollection represents the database collection of the NVM model.
var NVMCollection *vmdb.Collection

// AdressesCollection represents the database colltection of the Address model.
var AddressesCollection *vmdb.Collection

// CrewsCollection represents the database collection of the Crew model.
var CrewsCollection *vmdb.Collection

// ProfilesCollection represents the database collection of the Profile model.
var ProfilesCollection *vmdb.Collection

// AvatarCollection represents the database collection of the Avatar model.
var AvatarCollection *vmdb.Collection

// PoolRoleCollection represents the database collection of the PoolRole Collection.
var PoolRoleCollection *vmdb.Collection

var MailboxCollection *vmdb.Collection
var MessageCollection *vmdb.Collection

var ArtistCollection *vmdb.Collection
var ParticipationCollection *vmdb.Collection
var OrganizerCollection *vmdb.Collection
var EventCollection *vmdb.Collection
var SourceCollection *vmdb.Collection
var TakingCollection *vmdb.Collection

func InitialDatabase() {
	Database = vmdb.NewDatabase("pool-user").Connect()

	// UserCollection represents the database collection of the User model.
	UserCollection = Database.Collection(models.UserCollection).CreateIndex("email", true)

	// UserCrewCollection represents the database collection of the UserCrew model.
	UserCrewCollection = Database.Collection(models.UserCrewCollection).CreateIndex("user_id", true)

	// ActiveCollection represents the database collection of the Active model.
	ActiveCollection = Database.Collection(models.ActiveCollection).CreateIndex("user_id", true)

	// NVMCollection represents the database collection of the NVM model.
	NVMCollection = Database.Collection(models.NVMCollection).CreateIndex("user_id", true)

	// AdressesCollection represents the database colltection of the Address model.
	AddressesCollection = Database.Collection(models.AddressesCollection).CreateIndex("user_id", true)

	// CrewsCollection represents the database collection of the Crew model.
	CrewsCollection = Database.Collection(models.CrewCollection).CreateIndex("name", true)

	// ProfilesCollection represents the database collection of the Profile model.
	ProfilesCollection = Database.Collection(models.ProfilesCollection).CreateIndex("user_id", true)

	// AvatarCollection represents the database collection of the Avatar model.
	AvatarCollection = Database.Collection(models.AvatarCollection).CreateIndex("user_id", true)

	// PoolRoleCollection represents the database collection of the PoolRole Collection.
	PoolRoleCollection = Database.Collection(models.PoolRoleCollection).CreateMultiIndex(bson.D{{Key: "name", Value: 1}, {Key: "user_id", Value: 1}}, true)

	//
	MailboxCollection = Database.Collection(models.MailboxCollection)

	MessageCollection = Database.Collection(models.MessageCollection)
	ArtistCollection = Database.Collection("artists").CreateIndex("name", true)
	ParticipationCollection = Database.Collection("participations").CreateMultiIndex(
		bson.D{
			{Key: "user_id", Value: 1},
			{Key: "event_id", Value: 1},
		}, true)
	OrganizerCollection = Database.Collection("organizers").CreateIndex("name", true)
	EventCollection = Database.Collection("events")
	SourceCollection = Database.Collection("sources")
	TakingCollection = Database.Collection("takings").CreateIndex("event_id", true)
}

func ReloadDatabase() {
	userList := new([]models.User)
	var err error
	if err = UserCollection.Find(context.Background(), bson.D{}, userList); err != nil {
		log.Print(err)
		err = nil
	}
	for i := range *userList {
		vcago.Nats.Publish("pool.user.created", (*userList)[i])
	}
	crewList := new([]models.Crew)
	if err = CrewsCollection.Find(context.Background(), bson.D{}, crewList); err != nil {
		log.Print(err)
	}
	for i := range *crewList {
		vcago.Nats.Publish("pool.crew.created", (*crewList)[i])
	}
}
