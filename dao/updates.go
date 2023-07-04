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
	if !CheckUpdated(ctx, "update_usercrew3_mailbox") {
		UpdateUserCrewMaibox(ctx)
		InsertUpdate(ctx, "update_usercrew3_mailbox")
	}
	if !CheckUpdated(ctx, "update_delete_confirmed") {
		UpdateDeleteUnconfirmd(ctx)
		InsertUpdate(ctx, "update_delete_confirmed")
	}
	if !CheckUpdated(ctx, "update_confirm_admin") {
		UpdateConfirmAdmin(ctx)
		InsertUpdate(ctx, "update_confirm_admin")
	}
	if !CheckUpdated(ctx, "taking_currency1") {
		UpdateTakingCurrency(ctx)
		InsertUpdate(ctx, "taking_currency1")
	}
	if !CheckUpdated(ctx, "deposit_currency") {
		UpdateDepositCurrency(ctx)
		InsertUpdate(ctx, "deposit_currency")
	}
	if !CheckUpdated(ctx, "deposit_unit_currency") {
		UpdateDepositUnitCurrency(ctx)
		InsertUpdate(ctx, "deposit_unit_currency")
	}
	if !CheckUpdated(ctx, "taking_fix_no_income") {
		UpdateTakingFixNoIncome(ctx)
		InsertUpdate(ctx, "taking_fix_no_income")
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
		filter := bson.D{{Key: "crew_id", Value: (*crews)[i].ID}}
		update := bson.D{{Key: "mailbox_id", Value: (*crews)[i].MailboxID}}
		if _, err := UserCrewCollection.Collection.UpdateMany(ctx, filter, vmdb.UpdateSet(update)); err != nil {
			log.Print(err)
		}
	}
}

func UpdateDeleteUnconfirmd(ctx context.Context) {
	users := new([]models.User)
	filter := bson.D{{Key: "confirmed", Value: false}}
	if err := UserCollection.Find(ctx, filter, users); err != nil {
		log.Print(err)
	}
	for _, user := range *users {
		if err := UserDelete(ctx, user.ID); err != nil {
			log.Print(err)
		}
	}
}

func UpdateConfirmAdmin(ctx context.Context) {
	update := bson.D{{Key: "confirmed", Value: true}}
	filter := bson.D{{Key: "email", Value: "it@vivaconagua.org"}}
	if err := UserCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(update), nil); err != nil {
		log.Print(err)
	}
}

func UpdateTakingCurrency(ctx context.Context) {
	update := bson.D{{Key: "currency", Value: "EUR"}}
	if _, err := TakingCollection.Collection.UpdateMany(ctx, bson.D{}, vmdb.UpdateSet(update)); err != nil {
		return
	}
}

func UpdateDepositCurrency(ctx context.Context) {
	update := bson.D{{Key: "money.currency", Value: "EUR"}}
	if _, err := DepositCollection.Collection.UpdateMany(ctx, bson.D{}, vmdb.UpdateSet(update)); err != nil {
		return
	}
}

func UpdateDepositUnitCurrency(ctx context.Context) {
	update := bson.D{{Key: "money.currency", Value: "EUR"}}
	if _, err := DepositUnitCollection.Collection.UpdateMany(ctx, bson.D{}, vmdb.UpdateSet(update)); err != nil {
		return
	}
}

func UpdateTakingFixNoIncome(ctx context.Context) {
	update := bson.D{{Key: "state.no_income", Value: false}}
	if _, err := TakingCollection.Collection.UpdateMany(
		ctx,
		bson.D{{Key: "$not", Value: bson.D{{Key: "state.no_income", Value: true}}}},
		vmdb.UpdateSet(update),
	); err != nil {
		return
	}
}
