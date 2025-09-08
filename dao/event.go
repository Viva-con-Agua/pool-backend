package dao

import (
	"context"
	"pool-backend/models"
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
)

func EventInsert(ctx context.Context, i *models.EventCreate, token *models.AccessToken) (result *models.Event, err error) {
	if err = models.EventPermission(token); err != nil {
		return
	}
	event := i.EventDatabase(token)
	if !token.Roles.Validate("admin;employee;pool_employee") {
		crew := new(models.Crew)
		if crew, err = CrewGetByID(ctx, &models.CrewParam{ID: i.CrewID}, token); err != nil {
			return
		}
		event.OrganisationID = crew.OrganisationID
	}
	taking := event.TakingDatabase()
	event.TakingID = taking.ID
	if err = EventCollection.InsertOne(ctx, event); err != nil {
		return
	}
	if err = TakingCollection.InsertOne(ctx, taking); err != nil {
		return
	}
	eventActivity := models.NewActivity(event.CreatorID, "event", event.ID, "Event created", "created")
	if err = ActivityCollection.InsertOne(ctx, eventActivity); err != nil {
		return
	}
	takingActivity := models.NewActivity(event.CreatorID, "taking", taking.ID, "Taking automatically created for Event", "auto_created")
	if err = ActivityCollection.InsertOne(ctx, takingActivity); err != nil {
		return
	}
	filter := event.Match()
	if err = EventCollection.AggregateOne(ctx, models.EventPipeline(token).Match(filter).Pipe, &result); err != nil {
		return
	}
	history := result.NewEventStateHistory("", result.EventState.State, token)
	if err = EventStateHistoryInsert(ctx, history, token); err != nil {
		return
	}
	return
}

func EventGet(i *models.EventQuery, token *models.AccessToken) (result *[]models.ListEvent, list_size int64, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := i.PermittedFilter(token)
	sort := i.Sort()
	pipeline := models.EventPipeline(token).SortFields(sort).Match(filter).Sort(sort).Skip(i.Skip, 0).Limit(i.Limit, 100).Pipe
	result = new([]models.ListEvent)
	if err = EventCollection.Aggregate(ctx, pipeline, result); err != nil {
		return
	}

	count := vmod.Count{}
	var cErr error
	if cErr = EventCollection.AggregateOne(ctx, models.EventPipeline(token).Match(filter).Count().Pipe, &count); cErr != nil {
		list_size = 1
	} else {
		list_size = int64(count.Total)
	}
	return
}

func EventGetByID(ctx context.Context, i *models.EventParam, token *models.AccessToken) (result *models.Event, err error) {
	filter := i.PermittedFilter(token)
	if err = EventCollection.AggregateOne(
		ctx,
		models.EventPipeline(token).Match(filter).Pipe,
		&result,
	); err != nil {
		return
	}
	return
}

func EventGetInternalByID(ctx context.Context, i *models.EventParam) (result *models.Event, err error) {
	filter := i.FilterID()
	if err = EventCollection.AggregateOne(
		ctx,
		models.EventPipelinePublic().Match(filter).Pipe,
		&result,
	); err != nil {
		return
	}
	return
}

func EventAspGetByID(ctx context.Context, i *models.EventParam, token *models.AccessToken) (result *models.Event, err error) {
	filter := i.PermittedFilter(token)
	if err = EventCollection.AggregateOne(
		ctx,
		models.EventPipeline(token).Match(filter).Pipe,
		&result,
	); err != nil {
		return
	}
	return
}

func EventViewGetByID(ctx context.Context, i *models.EventParam) (result *models.EventPublic, err error) {
	filter := i.PublicFilter()
	if err = EventCollection.AggregateOne(
		ctx,
		models.EventPipelinePublic().Match(filter).Pipe,
		&result,
	); err != nil {
		return
	}
	return
}

func EventGetPublic(ctx context.Context, i *models.EventQuery) (result *[]models.EventPublic, list_size int64, err error) {
	filter := i.PublicFilter()
	sort := i.Sort()
	pipeline := models.EventPipelinePublic().SortFields(sort).Match(filter).Sort(sort).Skip(i.Skip, 0).Limit(i.Limit, 100).Pipe
	result = new([]models.EventPublic)
	if err = EventCollection.Aggregate(ctx, pipeline, result); err != nil {
		return
	}
	count := vmod.Count{}
	var cErr error
	cTx := context.Background()
	if cErr = EventCollection.AggregateOne(cTx, models.EventPipelinePublic().Match(filter).Count().Pipe, &count); cErr != nil {
		print(cErr)
		list_size = 1
	} else {
		list_size = int64(count.Total)
	}
	return
}

func EventsGetReceiverEvents(ctx context.Context, i *models.EventQuery, token *models.AccessToken) (result *[]models.EventPublic, err error) {
	filter := i.FilterEmailEvents(token)
	result = new([]models.EventPublic)
	if err = EventCollection.Aggregate(ctx, models.EventPipelinePublic().Match(filter).Limit(100, 100).Pipe, result); err != nil {
		return
	}
	return
}

func EventGetAps(ctx context.Context, i *models.EventQuery, token *models.AccessToken) (result *[]models.ListDetailsEvent, list_size int64, err error) {
	filter := i.FilterAsp(token)
	result = new([]models.ListDetailsEvent)

	sort := i.Sort()
	pipeline := models.EventPipelinePublic().SortFields(sort).Match(filter).Sort(sort).Skip(i.Skip, 0).Limit(i.Limit, 100).Pipe
	if err = EventCollection.Aggregate(ctx, pipeline, result); err != nil {
		return
	}
	count := vmod.Count{}
	var cErr error
	cTx := context.Background()
	if cErr = EventCollection.AggregateOne(cTx, models.EventPipelinePublic().Match(filter).Count().Pipe, &count); cErr != nil {
		print(cErr)
		list_size = 1
	} else {
		list_size = int64(count.Total)
	}
	return
}

func EventUpdate(ctx context.Context, i *models.EventUpdate, token *models.AccessToken) (result *models.Event, err error) {
	event := new(models.EventValidate)
	filter := i.PermittedFilter(token)
	if err = EventCollection.AggregateOne(ctx, models.EventPipelinePublic().Match(filter).Pipe, event); err != nil {
		return
	}
	if !token.Roles.Validate("admin;employee;pool_employee") {
		crew := new(models.Crew)
		if crew, err = CrewGetByID(ctx, &models.CrewParam{ID: token.CrewID}, token); err != nil {
			return
		}
		i.OrganisationID = crew.OrganisationID
	}
	taking := new(models.Taking)
	if taking, err = TakingGetByIDSystem(ctx, event.TakingID); err != nil {
		if !vmdb.ErrNoDocuments(err) {
			return
		}
	}
	event.Taking = *taking
	if err = i.EventStateValidation(token, event); err != nil {
		return
	}
	if err = EventCollection.UpdateOneAggregate(
		ctx,
		filter,
		vmdb.UpdateSet(i),
		&result,
		models.EventPipeline(token).Match(i.Match()).Pipe,
	); err != nil {
		return
	}
	if event.EventState.State != result.EventState.State {
		history := result.NewEventStateHistory(event.EventState.State, result.EventState.State, token)
		if err = EventStateHistoryInsert(ctx, history, token); err != nil {
			return
		}
		if result.EventState.State == "canceled" {
			EventParticipantsNotification(ctx, result, "event_cancel")
			updateTaking := bson.D{{Key: "state.no_income", Value: true}}
			filterTaking := bson.D{{Key: "_id", Value: event.Taking.ID}}
			if err = TakingCollection.UpdateOne(ctx, filterTaking, vmdb.UpdateSet(updateTaking), nil); err != nil {
				return
			}
		}
		if result.EventState.State == "published" ||
			result.EventState.State == "canceled" ||
			(result.EventState.State == "requested" && result.EventState.CrewConfirmation == "") {
			EventStateNotificationCreator(result)
			EventStateNotification(ctx, result)
		}
	} else if event.StartAt != result.StartAt ||
		event.EndAt != result.EndAt ||
		event.Location.PlaceID != result.Location.PlaceID ||
		event.EventASPID != result.EventASPID {
		EventParticipantsNotification(ctx, result, "event_update")
	}
	if event.EventASPID != result.EventASPID && result.EventASPID != token.ID {
		EventASPNotification(ctx, result)
	}
	if event.EndAt != i.EndAt {
		updateTaking := bson.D{{Key: "date_of_taking", Value: i.EndAt}}
		filterTaking := bson.D{{Key: "_id", Value: event.TakingID}}
		if err = TakingCollection.UpdateOne(ctx, filterTaking, vmdb.UpdateSet(updateTaking), nil); err != nil {
			return
		}
	}
	return
}

func EventDelete(ctx context.Context, i *models.EventParam, token *models.AccessToken) (err error) {
	if err = models.EventDeletePermission(token); err != nil {
		return
	}
	event := new(models.Event)
	filter := i.PermittedDeleteFilter(token)
	if err = EventCollection.FindOne(
		ctx,
		filter,
		event,
	); err != nil {
		return
	}
	deposit_unit := new(models.DepositUnit)
	if err = DepositUnitCollection.FindOne(
		ctx,
		bson.D{{Key: "taking_id", Value: event.TakingID}},
		deposit_unit,
	); err != nil {
		log.Info("No deposit units found")
	}
	if err = ParticipationCollection.TryDeleteMany(ctx, bson.D{{Key: "event_id", Value: i.ID}}); err != nil {
		return
	}
	if err = TakingCollection.TryDeleteMany(ctx, bson.D{{Key: "_id", Value: event.TakingID}}); err != nil {
		return
	}
	DepositCollection.TryDeleteMany(ctx, bson.D{{Key: "_id", Value: deposit_unit.DepositID}})
	DepositUnitCollection.TryDeleteMany(ctx, bson.D{{Key: "taking_id", Value: event.TakingID}})
	if err = EventCollection.DeleteOne(ctx, filter); err != nil {
		return
	}
	return
}

func EventImport(ctx context.Context, i *models.EventImport) (result *models.Event, err error) {
	event := i.EventDatabase()
	taking := event.TakingDatabase()
	event.TakingID = taking.ID

	admin := new(models.UserDatabase)
	adminFilter := bson.D{{Key: "email", Value: "f.wittmann@vivaconagua.org"}}
	if err = UserCollection.FindOne(
		ctx,
		adminFilter,
		admin,
	); err != nil {
		return
	}
	event.CreatorID = admin.ID
	event.InternalASPID = admin.ID
	event.EventState.InternalConfirmation = admin.ID
	/*
		if event.CrewID != "" {
			aspRole := new(models.RoleDatabase)
			if err = UserCrewCollection.AggregateOne(ctx, models.EventRolePipeline().Match(bson.D{{Key: "crew_id", Value: event.CrewID}}).Pipe, aspRole); err != nil {
				return
			}
			event.EventASPID = aspRole.UserID

			crew := new(models.Crew)
			if err = CrewsCollection.FindOne(ctx, bson.D{{Key: "_id", Value: event.CrewID}}, &crew); err != nil {
				return
			}
			event.OrganisationID = crew.OrganisationID

		} else {
			event.EventASPID = admin.ID
		}*/

	event.EventASPID = admin.ID
	if err = EventCollection.InsertOne(ctx, i); err != nil {
		return
	}
	if err = TakingCollection.InsertOne(ctx, taking); err != nil {
		return
	}
	eventActivity := models.NewActivity(event.CreatorID, "event", event.ID, "Event created", "created")
	if err = ActivityCollection.InsertOne(ctx, eventActivity); err != nil {
		return
	}
	takingActivity := models.NewActivity(event.CreatorID, "taking", taking.ID, "Taking automatically created for Event", "auto_created")
	if err = ActivityCollection.InsertOne(ctx, takingActivity); err != nil {
		return
	}

	if err = EventCollection.AggregateOne(ctx, models.EventImportPipeline().Match(event.Match()).Pipe, &result); err != nil {
		return
	}

	// Add participations
	for _, drops_user := range i.Participations {
		participant := new(models.UserDatabase)
		participantFilter := bson.D{{Key: "drops_id", Value: drops_user.DropsID}}
		if err = UserCollection.FindOne(
			ctx,
			participantFilter,
			participant,
		); err != nil {
			return
		}

		participation := drops_user.ParticipationDatabase()
		participation.EventID = result.ID
		participation.UserID = participant.ID
		participation.CrewID = result.CrewID
		if err = ParticipationCollection.InsertOne(ctx, participation); err != nil {
			return
		}
	}

	return
}

func EventSync(ctx context.Context, i *models.EventParam, token *models.AccessToken) (result *models.Event, err error) {
	if err = EventCollection.AggregateOne(ctx, models.EventPipeline(token).Match(i.Match()).Pipe, &result); err != nil {
		return
	}
	result.EditorID = result.Creator.ID
	go func() {
		if err = IDjango.Post(result, "/v1/pool/event/sync/"); err != nil {
			log.Print(err)
		}
	}()
	return
}

func EventParticipantsNotification(ctx context.Context, i *models.Event, template string) (err error) {
	filter := i.FilterParticipants()

	participants := new([]models.Participation)
	if err = ParticipationCollection.Aggregate(
		ctx,
		models.ParticipationPipeline().Match(filter).Pipe,
		participants,
	); err != nil {
		return
	}

	for _, participant := range *participants {
		mail := vcago.NewMailData(participant.User.Email, "pool-backend", template, "pool", participant.User.Country)
		mail.AddUser(participant.User.User())
		mail.AddContent(participant.ToContent())
		vcago.Nats.Publish("system.mail.job", mail)
		//notification := vcago.NewMNotificationData(participant.User.Email, "pool-backend", template, "pool", participant.User.Country, token.ID)
		//notification.AddUser(participant.User.User())
		//notification.AddContent(participant.ToContent())
		//vcago.Nats.Publish("system.notification.job", notification)
	}

	return
}

func EventASPNotification(ctx context.Context, i *models.Event) (err error) {

	if i.EventASPID == "" {
		return vcago.NewNotFound(models.EventCollection, i)
	}

	user := new(models.User)
	if user, err = UsersGetByID(ctx, &models.UserParam{ID: i.EventASPID}); err != nil {
		return
	}

	template := "event_asp"
	mail := vcago.NewMailData(user.Email, "pool-backend", template, "pool", user.Country)
	mail.AddUser(user.User())
	mail.AddContent(i.ToContent())
	vcago.Nats.Publish("system.mail.job", mail)

	//notification := vcago.NewMNotificationData(user.Email, "pool-backend", template, user.Country, token.ID)
	//notification.AddUser(user.User())
	//notification.AddContent(i.ToContent())
	//vcago.Nats.Publish("system.notification.job", notification)
	return
}

func EventStateNotification(ctx context.Context, i *models.Event) (err error) {

	if i.EventASPID == "" {
		return vcago.NewNotFound(models.EventCollection, i)
	}

	eventAsp := new(models.User)
	if eventAsp, err = UsersGetByID(ctx, &models.UserParam{ID: i.EventASPID}); err != nil {
		return
	}

	notifyAboutStateChange(eventAsp, i)
	return
}

func notifyAboutStateChange(notifyUser *models.User, i *models.Event) {
	template := "event_state"
	mail := vcago.NewMailData(notifyUser.Email, "pool-backend", template, "pool", notifyUser.Country)
	mail.AddUser(notifyUser.User())
	mail.AddContent(i.ToContent())
	vcago.Nats.Publish("system.mail.job", mail)

	//notification := vcago.NewMNotificationData(notifyUser.Email, "pool-backend", template, notifyUser.Country, token.ID)
	//notification.AddUser(notifyUser.User())
	//notification.AddContent(i.ToContent())
	//vcago.Nats.Publish("system.notification.job", notification)
}

func EventStateNotificationCreator(i *models.Event) (err error) {
	// if event_asp is creator -> stop.
	if i.EventASPID == i.Creator.ID {
		return
	}

	notifyAboutStateChange(&i.Creator, i)
	return
}

func EventHistoryAdminNotification(ctx context.Context, data []models.EventStateHistoryNotification) (err error) {
	mail := vcago.NewMailData("netzwerk@vivaconagua.org", "pool-backend", "events_published", "pool", "de")
	mail.AddContent(models.EventHistoryAdminContent(data))
	vcago.Nats.Publish("system.mail.job", mail)

	mail2 := vcago.NewMailData("festival@vivaconagua.org", "pool-backend", "events_published", "pool", "de")
	mail2.AddContent(models.EventHistoryAdminContent(data))
	vcago.Nats.Publish("system.mail.job", mail2)
	return
}

func EventHistoryCrewNotification(ctx context.Context, data_collection map[string][]models.EventStateHistoryNotification) (err error) {
	for email, data := range data_collection {
		mail := vcago.NewMailData(email, "pool-backend", "events_crew_published", "pool", "de")
		mail.AddContent(models.EventHistoryAdminContent(data))
		vcago.Nats.Publish("system.mail.job", mail)
	}
	return
}
