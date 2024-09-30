package dao

import (
	"context"
	"log"
	"pool-backend/models"
	"time"

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
	if !CheckUpdated(ctx, "taking_no_income_event_canceled") {
		UpdateEventCanceledNoIncome(ctx)
		InsertUpdate(ctx, "taking_no_income_event_canceled")
	}
	if !CheckUpdated(ctx, "currency_problem") {
		UpdateDepositCurrency(ctx)
		UpdateDepositUnitCurrency(ctx)
		UpdateTakingCurrency(ctx)
		InsertUpdate(ctx, "currency_problem")
	}
	if !CheckUpdated(ctx, "date_of_taking_1") {
		UpdateDateOfTaking1(ctx)
		InsertUpdate(ctx, "date_of_taking_1")
	}
	if !CheckUpdated(ctx, "birthdate_1") {
		UpdateProfileBirthdate(ctx)
		InsertUpdate(ctx, "birthdate_1")
	}
	if !CheckUpdated(ctx, "event_applications") {
		UpdateEventApplications(ctx)
		InsertUpdate(ctx, "event_applications")
	}
	if !CheckUpdated(ctx, "last_login_date_1") {
		UpdateSetLastLoginDate(ctx)
		InsertUpdate(ctx, "last_login_date_1")
	}
	if !CheckUpdated(ctx, "create_default_organisation") {
		CreateDefaultOrganisation(ctx)
		InsertUpdate(ctx, "create_default_organisation")
	}
	if !CheckUpdated(ctx, "update_deposit_units_1") {
		UpdateDepositUnitNorms(ctx)
		InsertUpdate(ctx, "update_deposit_units_1")
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
	filter := bson.D{{Key: "money.currency", Value: ""}}
	if _, err := TakingCollection.Collection.UpdateMany(ctx, filter, vmdb.UpdateSet(update)); err != nil {
		return
	}
}

func UpdateDepositCurrency(ctx context.Context) {
	update := bson.D{{Key: "money.currency", Value: "EUR"}}
	filter := bson.D{{Key: "money.currency", Value: ""}}
	if _, err := DepositCollection.Collection.UpdateMany(ctx, filter, vmdb.UpdateSet(update)); err != nil {
		return
	}
}

func UpdateDepositUnitCurrency(ctx context.Context) {
	update := bson.D{{Key: "money.currency", Value: "EUR"}}
	filter := bson.D{{Key: "money.currency", Value: ""}}
	if _, err := DepositUnitCollection.Collection.UpdateMany(ctx, filter, vmdb.UpdateSet(update)); err != nil {
		return
	}
}

func UpdateEventCanceledNoIncome(ctx context.Context) {
	filterEvent := bson.D{{Key: "event_state.state", Value: "canceled"}}
	eventResult := []models.Event{}
	if err := EventCollection.Find(ctx, filterEvent, &eventResult); err != nil {
		return
	}
	for _, value := range eventResult {
		updateTaking := bson.D{{Key: "state.no_income", Value: true}}
		filterTaking := bson.D{{Key: "_id", Value: value.TakingID}}
		if err := TakingCollection.UpdateOne(ctx, filterTaking, vmdb.UpdateSet(updateTaking), nil); err != nil {
			return
		}
	}
}

func UpdateDateOfTaking1(ctx context.Context) {
	eventList := []models.Event{}
	if err := EventCollection.Find(ctx, bson.D{{}}, &eventList); err != nil {
		log.Print(err)
	}
	for _, event := range eventList {
		update := bson.D{{Key: "date_of_taking", Value: event.EndAt}}
		filter := bson.D{{Key: "_id", Value: event.TakingID}}
		if err := TakingCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(update), nil); err != nil {
			log.Print(err)
		}
	}
}

func UpdateProfileBirthdate(ctx context.Context) {
	profileList := []models.ProfileUpdate{}
	if err := ProfileCollection.Find(ctx, bson.D{{}}, &profileList); err != nil {
		log.Print(err)
	}
	for _, profile := range profileList {
		birthdate := time.Unix(profile.Birthdate, 0)
		if profile.Birthdate != 0 {
			profile.BirthdateDatetime = birthdate.Format("2006-01-02")
		} else {
			profile.BirthdateDatetime = ""
		}
		filter := bson.D{{Key: "_id", Value: profile.ID}}
		if err := ProfileCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(profile), nil); err != nil {
			log.Print(err)
		}
	}
}

func UpdateEventApplications(ctx context.Context) {
	eventList := []models.Event{}
	if err := EventCollection.Aggregate(ctx, models.EventPipeline(&models.AccessToken{ID: ""}).Pipe, &eventList); err != nil {
		log.Print(err)
	}
	for _, event := range eventList {
		confirmed, rejected, requested, withdrawn, total := 0, 0, 0, 0, 0

		for _, p := range event.Participation {
			switch p.Status {
			case "confirmed":
				confirmed++
			case "rejected":
				rejected++
			case "requested":
				requested++
			case "withdrawn":
				withdrawn++
			}
			total++
		}
		update := bson.D{{Key: "applications", Value: models.EventApplications{
			Confirmed: confirmed, Rejected: rejected, Requested: requested, Withdrawn: withdrawn, Total: total,
		}}}
		filter := bson.D{{Key: "_id", Value: event.ID}}
		if err := EventCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(update), nil); err != nil {
			log.Print(err)
		}
	}
}
func UpdateSetLastLoginDate(ctx context.Context) {
	update := bson.D{{Key: "last_login_date", Value: time.Now().Unix()}}
	if err := UserCollection.UpdateMany(ctx, bson.D{{}}, vmdb.UpdateSet(update)); err != nil {
		log.Print(err)
	}
}

func CreateDefaultOrganisation(ctx context.Context) {
	i := models.OrganisationCreate{
		Name:         "Viva con Agua de Sankt Pauli e.V.",
		Abbreviation: "VcA DE",
		Email:        "pool@vivaconagua.org",
	}
	result := new(models.Organisation)
	result = i.Organisation()
	if err := OrganisationCollection.InsertOne(ctx, result); err != nil {
		log.Print(err)
	}
	update := bson.D{{Key: "organisation_id", Value: result.ID}}
	if err := CrewsCollection.UpdateMany(ctx, bson.D{}, vmdb.UpdateSet(update)); err != nil {
		log.Print(err)
	}
	if err := UserCrewCollection.UpdateMany(ctx, bson.D{}, vmdb.UpdateSet(update)); err != nil {
		log.Print(err)
	}
	if err := EventCollection.UpdateMany(ctx, bson.D{}, vmdb.UpdateSet(update)); err != nil {
		log.Print(err)
	}
	filter := vmdb.NewFilter()
	filter.ElemMatchList("system_roles", "name", []string{"employee", "pool_employee", "pool_finance"})
	if err := UserCollection.UpdateMany(ctx, filter.Bson(), vmdb.UpdateSet(update)); err != nil {
		log.Print(err)
	}
}

func UpdateDepositUnitNorms(ctx context.Context) {
	filterDonation := vmdb.NewFilter()
	filterDonation.EqualStringList("value", []string{"unknown", "can", "box", "gl", "other"})
	filterEco := vmdb.NewFilter()
	filterEco.EqualStringList("value", []string{"merch", "other_ec"})
	updateDonation := bson.D{{Key: "$set", Value: bson.D{{Key: "norms", Value: "donation"}}}}
	updateEco := bson.D{{Key: "$set", Value: bson.D{{Key: "norms", Value: "economic"}}}}
	if err := SourceCollection.UpdateMany(ctx, filterDonation.Bson(), updateDonation); err != nil {
		log.Print(err)
	}
	if err := SourceCollection.UpdateMany(ctx, filterEco.Bson(), updateEco); err != nil {
		log.Print(err)
	}
}
