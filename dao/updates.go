package dao

import (
	"context"
	"log"
	"pool-backend/models"
	"time"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateDatabase() {
	ctx := context.Background()
	if !Updates.Check(ctx, "update_crew_mailbox") {
		UpdateCrewMaibox(ctx)
		Updates.Insert(ctx, "update_crew_mailbox")
	}
	if !Updates.Check(ctx, "update_usercrew3_mailbox") {
		UpdateUserCrewMaibox(ctx)
		Updates.Insert(ctx, "update_usercrew3_mailbox")
	}
	if !Updates.Check(ctx, "update_delete_confirmed") {
		UpdateDeleteUnconfirmd(ctx)
		Updates.Insert(ctx, "update_delete_confirmed")
	}
	if !Updates.Check(ctx, "update_confirm_admin") {
		UpdateConfirmAdmin(ctx)
		Updates.Insert(ctx, "update_confirm_admin")
	}
	if !Updates.Check(ctx, "taking_currency1") {
		UpdateTakingCurrency(ctx)
		Updates.Insert(ctx, "taking_currency1")
	}
	if !Updates.Check(ctx, "deposit_currency") {
		UpdateDepositCurrency(ctx)
		Updates.Insert(ctx, "deposit_currency")
	}
	if !Updates.Check(ctx, "deposit_unit_currency") {
		UpdateDepositUnitCurrency(ctx)
		Updates.Insert(ctx, "deposit_unit_currency")
	}
	if !Updates.Check(ctx, "taking_no_income_event_canceled") {
		UpdateEventCanceledNoIncome(ctx)
		Updates.Insert(ctx, "taking_no_income_event_canceled")
	}
	if !Updates.Check(ctx, "currency_problem") {
		UpdateDepositCurrency(ctx)
		UpdateDepositUnitCurrency(ctx)
		UpdateTakingCurrency(ctx)
		Updates.Insert(ctx, "currency_problem")
	}
	if !Updates.Check(ctx, "date_of_taking_1") {
		UpdateDateOfTaking1(ctx)
		Updates.Insert(ctx, "date_of_taking_1")
	}
	if !Updates.Check(ctx, "birthdate_1") {
		UpdateProfileBirthdate(ctx)
		Updates.Insert(ctx, "birthdate_1")
	}
	if !Updates.Check(ctx, "event_applications") {
		UpdateEventApplications(ctx)
		Updates.Insert(ctx, "event_applications")
	}
	if !Updates.Check(ctx, "last_login_date_1") {
		UpdateSetLastLoginDate(ctx)
		Updates.Insert(ctx, "last_login_date_1")
	}
	if !Updates.Check(ctx, "create_default_organisation") {
		CreateDefaultOrganisation(ctx)
		Updates.Insert(ctx, "create_default_organisation")
	}
	if !Updates.Check(ctx, "update_deposit_units_4") {
		UpdateDepositUnitNorms(ctx)
		Updates.Insert(ctx, "update_deposit_units_4")
	}
	if !Updates.Check(ctx, "publish_roles_initial") {
		log.Print("publish_roles_initial")
		PublishRoles()
		Updates.Insert(ctx, "publish_roles_initial")
	}
	if !Updates.Check(ctx, "date_of_deposit") {
		UpdateDateOfDeposit(ctx)
		Updates.Insert(ctx, "date_of_deposit")
	}
	if !Updates.Check(ctx, "event_applications_fix_2") {
		log.Print("event_applications_fix")
		UpdateEventApplications(ctx)
		Updates.Insert(ctx, "event_applications_fix_2")
	}
	if !Updates.Check(ctx, "event_crew_id_fix") {
		log.Print("event_crew_id_fix")
		UpdateEventCrewIDs(ctx)
		Updates.Insert(ctx, "event_crew_id_fix")
	}
	if !Updates.Check(ctx, "organisation_options") {
		log.Print("organisation_options")
		UpdateOrganisationOptions(ctx)
		Updates.Insert(ctx, "organisation_options")
	}
	if !Updates.Check(ctx, "update_nvm_to_user") {
		log.Print("update_nvm_to_user")
		UpdateNVMtoUser(ctx)
		Updates.Insert(ctx, "update_nvm_to_user")
	}
	if !Updates.Check(ctx, "update_active_to_user") {
		log.Print("update_active_to_user")
		UpdateActivetoUser(ctx)
		Updates.Insert(ctx, "update_active_to_user")
	}
	if !Updates.Check(ctx, "update_avatar_to_user") {
		log.Print("update_avatar_to_user")
		UpdateAvatartoUser(ctx)
		Updates.Insert(ctx, "update_avatar_to_user")
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

func UpdateEventCrewIDs(ctx context.Context) {
	eventList := []models.EventUpdate{}
	eventFilter := bson.D{{Key: "crew_id", Value: ""}, {Key: "event_asp_id", Value: bson.D{{Key: "$ne", Value: ""}}}}

	if err := EventCollection.Find(ctx, eventFilter, &eventList); err != nil {
		log.Print(err)
	}
	for _, event := range eventList {
		user := models.UserCrew{}
		if err := UserCrewCollection.FindOne(ctx, bson.D{{Key: "user_id", Value: event.EventASPID}}, &user); err != nil {
			log.Print(err)
		}
		if event.EventASPID == event.InternalASPID {
			continue
		}
		filter := bson.D{{Key: "_id", Value: event.ID}}
		event.CrewID = user.CrewID
		if err := EventCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(event), nil); err != nil {
			log.Print(err)
		}
	}
}

func UpdateEventApplications(ctx context.Context) {
	eventList := []models.Event{}
	if err := EventCollection.Aggregate(ctx, models.EventPipeline(&models.AccessToken{ID: "internal"}).Pipe, &eventList); err != nil {
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

func UpdateDateOfDeposit(ctx context.Context) {
	filter := vmdb.NewFilter()
	filter.EqualStringList("status", []string{"wait", "confirmed"})
	deposits := []models.Deposit{}
	if err := DepositCollection.Find(ctx, bson.D{{}}, &deposits); err != nil {
		log.Print(err)
	}
	for _, entry := range deposits {
		updateFilter := bson.D{{Key: "_id", Value: entry.ID}}
		update := bson.D{{Key: "date_of_deposit", Value: entry.Modified.Created}}
		if err := DepositCollection.UpdateOne(ctx, updateFilter, vmdb.UpdateSet(update), nil); err != nil {
			log.Print(err)
		}
	}
}

func UpdateOrganisationOptions(ctx context.Context) {
	update := bson.D{{Key: "options", Value: []string{models.OptionActiv, models.OptionNVM, models.OptionVolunteerCert}}}
	if err := OrganisationCollection.UpdateMany(ctx, bson.D{}, vmdb.UpdateSet(update)); err != nil {
		log.Print(err)
	}
}

func UpdateNVMtoUser(ctx context.Context) {
	nvm := []models.NVM{}
	if err := NVMCollection.Find(ctx, bson.D{}, &nvm); err != nil {
		log.Print(err)
	}
	for _, entry := range nvm {
		filter := bson.D{{Key: "_id", Value: entry.UserID}}
		update := bson.D{{Key: "nvm", Value: entry}}
		if err := UserCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(update), nil); err != nil {
			log.Print(err)
		}
	}
}

func UpdateActivetoUser(ctx context.Context) {
	active := []models.Active{}
	if err := ActiveCollection.Find(ctx, bson.D{}, &active); err != nil {
		log.Print(err)
	}
	for _, entry := range active {
		filter := bson.D{{Key: "_id", Value: entry.UserID}}
		update := bson.D{{Key: "active", Value: entry}}
		if err := UserCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(update), nil); err != nil {
			log.Print(err)
		}
	}
}

func UpdateAvatartoUser(ctx context.Context) {
	active := []models.Active{}
	if err := AvatarCollection.Find(ctx, bson.D{}, &active); err != nil {
		log.Print(err)
	}
	for _, entry := range active {
		filter := bson.D{{Key: "_id", Value: entry.UserID}}
		update := bson.D{{Key: "avatar", Value: entry}}
		if err := UserCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(update), nil); err != nil {
			log.Print(err)
		}
	}
}
