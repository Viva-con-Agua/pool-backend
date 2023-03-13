package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Updated struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

func CheckUpdated(ctx context.Context, name string) bool {
	update := new(Updated)
	if err := UpdateCollection.FindOne(ctx, bson.D{{Key: "name", Value: name}}, update); err != nil {
		if vmdb.ErrNoDocuments(err) {
			return false
		}
		log.Print(err)
	}
	return true
}

func InsertUpdate(ctx context.Context, name string) {
	update := &Updated{ID: uuid.NewString(), Name: name}
	if err := UpdateCollection.InsertOne(ctx, update); err != nil {
		log.Print(err)
	}

}

func UpdateDatabase() {
	ctx := context.Background()
	if !CheckUpdated(ctx, "update_crew_mailbox") {
		UpdateCrewMaibox(ctx)
		InsertUpdate(ctx, "update_crew_mailbox")
	}
	if !CheckUpdated(ctx, "update_usercrew_mailbox") {
		UpdateCrewMaibox(ctx)
		InsertUpdate(ctx, "update_usercrew_mailbox")
	}
}

func UpdateCrewMaibox(ctx context.Context) {
	crews := new([]models.Crew)
	if err := CrewsCollection.Find(ctx, bson.D{{Key: "mailbox_id", Value: ""}}, crews); err != nil {
		log.Print(err)
	}
	for i := range *crews {
		mailbox := models.NewMailboxDatabase("crew")
		if err := MailboxCollection.InsertOne(ctx, mailbox); err != nil {
			log.Print()
		}
		filter := bson.D{{Key: "_id", Value: (*crews)[i].ID}}
		update := bson.D{{Key: "mailbox_id", Value: mailbox.ID}}
		if err := CrewsCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(update), nil); err != nil {
			log.Print(err)
		}
	}
}

func UpdateUserCrewMaibox(ctx context.Context) {
	crews := new([]models.Crew)
	if err := CrewsCollection.Find(ctx, bson.D{}, crews); err != nil {
		log.Print(err)
	}
	for i := range *crews {
		filter := bson.D{{Key: "crew_id", Value: (*crews)[i].ID}, {Key: "mailbox_id", Value: ""}}
		update := bson.D{{Key: "mailbox_id", Value: (*crews)[i].ID}}
		if _, err := UserCrewCollection.Collection.UpdateMany(ctx, filter, vmdb.UpdateSet(update)); err != nil {
			log.Print(err)
		}
	}
}
