package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
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

	// ProfileCollection represents the database collection of the Profile model.
	ProfileCollection *vmdb.Collection

	// AvatarCollection represents the database collection of the Avatar model.
	AvatarCollection *vmdb.Collection

	// PoolRoleCollection represents the database collection of the PoolRole Collection.
	PoolRoleCollection        *vmdb.Collection
	PoolRoleHistoryCollection *vmdb.Collection

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
	UserViewCollection         *vmdb.Collection
	EventViewCollection        *vmdb.Collection
	PublicEventViewCollection  *vmdb.Collection

	NewsletterCollection *vmdb.Collection

	DepositUnitTakingPipe  = vmdb.NewPipeline()
	ParticipationEventPipe = vmdb.NewPipeline()
	ActitityUserPipe       = vmdb.NewPipeline()
	UserPipe               = vmdb.NewPipeline()
	EventPipe              = vmdb.NewPipeline()
	PublicEventPipe        = vmdb.NewPipeline()
	UpdateCollection       *vmdb.Collection
	ReceiptFileCollection  *vmdb.Collection

	TestLogin bool
)

func InitialDatabase() {
	Database = vmdb.NewDatabase("pool-backend").Connect()

	// UserCollection represents the database collection of the User model.
	UserCollection = Database.Collection(models.UserCollection).CreateIndex("email", true)

	// UserCrewCollection represents the database collection of the UserCrew model.
	UserCrewCollection = Database.Collection(models.UserCrewCollection).CreateIndex("user_id", true).CreateIndex("crew_id", false).CreateMultiIndex(bson.D{{Key: "crew_id", Value: 1}, {Key: "user_id", Value: 1}}, true)

	// ActiveCollection represents the database collection of the Active model.
	ActiveCollection = Database.Collection(models.ActiveCollection).CreateIndex("user_id", true).CreateIndex("crew_id", false).CreateMultiIndex(bson.D{{Key: "crew_id", Value: 1}, {Key: "user_id", Value: 1}}, true)

	// NVMCollection represents the database collection of the NVM model.
	NVMCollection = Database.Collection(models.NVMCollection).CreateIndex("user_id", true).CreateIndex("crew_id", false).CreateMultiIndex(bson.D{{Key: "crew_id", Value: 1}, {Key: "user_id", Value: 1}}, true)

	// AdressesCollection represents the database colltection of the Address model.
	AddressesCollection = Database.Collection(models.AddressesCollection).CreateIndex("user_id", true)

	// CrewsCollection represents the database collection of the Crew model.
	CrewsCollection = Database.Collection(models.CrewCollection).CreateIndex("name", true)

	// ProfileCollection represents the database collection of the Profile model.
	ProfileCollection = Database.Collection(models.ProfileCollection).CreateIndex("user_id", true)

	// AvatarCollection represents the database collection of the Avatar model.
	AvatarCollection = Database.Collection(models.AvatarCollection).CreateIndex("user_id", true)

	// PoolRoleCollection represents the database collection of the PoolRole Collection.
	PoolRoleCollection = Database.Collection(models.PoolRoleCollection).CreateIndex("user_id", false).CreateMultiIndex(bson.D{{Key: "name", Value: 1}, {Key: "user_id", Value: 1}}, true)
	PoolRoleHistoryCollection = Database.Collection(models.PoolRoleHistoryCollection).CreateIndex("user_id", false).CreateIndex("crew_id", false)

	//
	MailboxCollection = Database.Collection(models.MailboxCollection)

	MessageCollection = Database.Collection(models.MessageCollection).CreateIndex("user_id", false).CreateIndex("mailbox_id", false)
	ArtistCollection = Database.Collection(models.ArtistCollection).CreateIndex("name", true)
	ParticipationCollection = Database.Collection(models.ParticipationCollection).CreateIndex("user_id", false).CreateMultiIndex(
		bson.D{
			{Key: "user_id", Value: 1},
			{Key: "event_id", Value: 1},
		}, true)
	OrganizerCollection = Database.Collection(models.OrganizerCollection).CreateIndex("name", true)
	EventCollection = Database.Collection(models.EventCollection)
	SourceCollection = Database.Collection(models.SourceCollection)
	TakingCollection = Database.Collection(models.TakingCollection).CreateIndex("crew_id", false)
	DepositCollection = Database.Collection(models.DepositCollection)
	DepositUnitCollection = Database.Collection(models.DepositUnitCollection).CreateMultiIndex(bson.D{{Key: "taking_id", Value: 1}, {Key: "deposit_id", Value: 1}}, true).CreateIndex("taking_id", false).CreateIndex("deposit_id", false)

	FSChunkCollection = Database.Collection(models.FSChunkCollection)
	FSFilesCollection = Database.Collection(models.FSFilesCollection)
	ActivityCollection = Database.Collection(models.ActivityCollection)

	NewsletterCollection = Database.Collection(models.NewsletterCollection).CreateIndex("user_id", false).CreateMultiIndex(
		bson.D{{Key: "user_id", Value: 1}, {Key: "value", Value: 1}},
		true,
	)

	ReasonForPaymentCollection = Database.Collection(models.ReasonForPaymentCollection)

	UserPipe.LookupUnwind(models.AddressesCollection, "_id", "user_id", "address")
	UserPipe.LookupUnwind(models.ProfileCollection, "_id", "user_id", "profile")
	UserPipe.LookupUnwind(models.UserCrewCollection, "_id", "user_id", "crew")
	UserPipe.LookupUnwind(models.ActiveCollection, "_id", "user_id", "active")
	UserPipe.LookupUnwind(models.NVMCollection, "_id", "user_id", "nvm")
	UserPipe.Lookup(models.PoolRoleCollection, "_id", "user_id", "pool_roles")
	Database.Database.CreateView(
		context.Background(),
		models.UserView,
		models.UserCollection,
		UserPipe.Pipe,
	)
	UserViewCollection = Database.Collection(models.UserView)

	PublicEventPipe.Lookup(models.ParticipationCollection, "_id", "event_id", "participations")
	PublicEventPipe.LookupUnwind(models.OrganizerCollection, "organizer_id", "_id", "organizer")
	PublicEventPipe.LookupList(models.ArtistCollection, "artist_ids", "_id", "artists")
	PublicEventPipe.LookupUnwind(models.CrewCollection, "crew_id", "_id", "crew")
	Database.Database.CreateView(
		context.Background(),
		models.PublicEventView,
		models.EventCollection,
		models.EventPipelinePublic().Pipe,
	)
	PublicEventViewCollection = Database.Collection(models.PublicEventView)

	EventPipe.Lookup(models.ParticipationCollection, "_id", "event_id", "participations")
	EventPipe.LookupUnwind(models.OrganizerCollection, "organizer_id", "_id", "organizer")
	EventPipe.LookupList(models.ArtistCollection, "artist_ids", "_id", "artists")
	EventPipe.LookupUnwind(models.CrewCollection, "crew_id", "_id", "crew")
	Database.Database.CreateView(
		context.Background(),
		models.EventView,
		models.EventCollection,
		models.EventPipeline(&vcapool.AccessToken{ID: ""}),
	)
	EventViewCollection = Database.Collection(models.EventView)

	DepositUnitTakingPipe.LookupUnwind(models.DepositCollection, "deposit_id", "_id", "deposit")
	Database.Database.CreateView(
		context.Background(),
		models.DepositUnitTakingView,
		models.DepositUnitCollection,
		DepositUnitTakingPipe.Pipe,
	)

	ParticipationEventPipe.LookupUnwind(models.EventCollection, "event_id", "_id", "event")
	Database.Database.CreateView(
		context.Background(),
		models.ParticipationEventView,
		models.ParticipationCollection,
		ParticipationEventPipe.Pipe,
	)
	ActitityUserPipe.LookupUnwind(models.UserCollection, "user_id", "_id", "user")
	Database.Database.CreateView(
		context.Background(),
		models.ActivityUserView,
		models.ActivityCollection,
		ActitityUserPipe.Pipe,
	)
	UpdateCollection = Database.Collection("updates").CreateIndex("name", true)
	ReceiptFileCollection = Database.Collection(models.ReceiptFileCollection).CreateIndex("deposit_id", false)
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

var IDjango = new(vcago.IDjangoHandler)

func InitialIDjango() {
	IDjango.URL = vcago.Settings.String("IDJANGO_URL", "n", "https://idjango.dev.vivaconagua.org")
	IDjango.Key = vcago.Settings.String("IDJANGO_KEY", "n", "")
	IDjango.Export = vcago.Settings.Bool("IDJANGO_EXPORT", "n", false)
}

func InitialTestLogin() {
	TestLogin = vcago.Settings.Bool("API_TEST_LOGIN", "n", false)
}
