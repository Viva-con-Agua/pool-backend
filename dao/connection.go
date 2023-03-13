package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	Database *vmdb.Database

	// UserCollection represents the database collection of the User model.
	UserCollection *vmdb.Collection

	// UserCrewCollection represents the database collection of the UserCrew model.
	UserCrewCollection *vmdb.Collection

	// ActiveCollection represents the database collection of the Active model.
	ActiveCollection *vmdb.Collection

	// NVMCollection represents the database collection of the NVM model.
	NVMCollection *vmdb.Collection

	// AdressesCollection represents the database colltection of the Address model.
	AddressesCollection *vmdb.Collection

	// CrewsCollection represents the database collection of the Crew model.
	CrewsCollection *vmdb.Collection

	// ProfilesCollection represents the database collection of the Profile model.
	ProfilesCollection *vmdb.Collection

	// AvatarCollection represents the database collection of the Avatar model.
	AvatarCollection *vmdb.Collection

	// PoolRoleCollection represents the database collection of the PoolRole Collection.
	PoolRoleCollection *vmdb.Collection

	MailboxCollection *vmdb.Collection
	MessageCollection *vmdb.Collection

	ArtistCollection        *vmdb.Collection
	ParticipationCollection *vmdb.Collection
	OrganizerCollection     *vmdb.Collection
	EventCollection         *vmdb.Collection

	SourceCollection      *vmdb.Collection
	TakingCollection      *vmdb.Collection
	DepositCollection     *vmdb.Collection
	DepositUnitCollection *vmdb.Collection

	FSChunkCollection *vmdb.Collection
	FSFilesCollection *vmdb.Collection

	ActivityCollection *vmdb.Collection

	ReasonForPaymentCollection *vmdb.Collection

	NewsletterCollection *vmdb.Collection

	DepositUnitTakingPipe  = vmdb.NewPipeline()
	ParticipationEventPipe = vmdb.NewPipeline()
	ActitityUserPipe       = vmdb.NewPipeline()
	UpdateCollection       *vmdb.Collection
)

func InitialDatabase() {
	Database = vmdb.NewDatabase("pool-backend").Connect()

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
	TakingCollection = Database.Collection("takings")
	DepositCollection = Database.Collection("deposits")
	DepositUnitCollection = Database.Collection("deposit_units").CreateMultiIndex(bson.D{{Key: "taking_id", Value: 1}, {Key: "deposit_id", Value: 1}}, true)

	FSChunkCollection = Database.Collection("fs.chunks")
	FSFilesCollection = Database.Collection("fs.files")
	ActivityCollection = Database.Collection("activities")

	NewsletterCollection = Database.Collection("newsletters").CreateMultiIndex(
		bson.D{{Key: "user_id", Value: 1}, {Key: "value", Value: 1}},
		true,
	)

	ReasonForPaymentCollection = Database.Collection("reason_for_payment")
	DepositUnitTakingPipe.LookupUnwind("deposits", "deposit_id", "_id", "deposit")
	Database.Database.CreateView(
		context.Background(),
		"deposit_unit_taking",
		"deposit_units",
		DepositUnitTakingPipe.Pipe,
	)

	ParticipationEventPipe.LookupUnwind("events", "event_id", "_id", "event")
	Database.Database.CreateView(
		context.Background(),
		"participations_event",
		"participations",
		ParticipationEventPipe.Pipe,
	)
	ActitityUserPipe.LookupUnwind("users", "user_id", "_id", "user")
	Database.Database.CreateView(
		context.Background(),
		"activity_user",
		"activities",
		ActitityUserPipe.Pipe,
	)
	UpdateCollection = Database.Collection("updates").CreateIndex("name", true)
}

func FixDatabase() {
	eventList := new([]models.EventDatabase)
	var err error
	if err = EventCollection.Find(
		context.Background(),
		bson.D{{Key: "$or", Value: bson.A{
			bson.D{{Key: "taking_id", Value: ""}},
			bson.D{{Key: "taking_id", Value: bson.D{{Key: "$exists", Value: false}}}},
		}}},
		eventList,
	); err != nil {
		log.Print(err)
	}
	for i := range *eventList {
		taking := models.TakingDatabase{
			ID:       uuid.NewString(),
			Modified: vmod.NewModified(),
		}
		event := (*eventList)[i]
		event.TakingID = taking.ID
		event.Modified.Update()
		if err = EventCollection.UpdateOne(
			context.Background(),
			bson.D{{Key: "_id", Value: event.ID}},
			bson.D{{Key: "$set", Value: event}},
			nil,
		); err != nil {
			log.Print(err)
		}
		if err = TakingCollection.InsertOne(context.Background(), taking); err != nil {
			log.Print(err)
		}
	}
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
