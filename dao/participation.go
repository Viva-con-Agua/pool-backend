package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"go.mongodb.org/mongo-driver/bson"
)

func ParticipationInsert(ctx context.Context, i *models.ParticipationCreate, token *models.AccessToken) (result *models.Participation, err error) {

	event := new(models.Event)
	if err = EventCollection.FindOne(
		ctx,
		bson.D{{Key: "_id", Value: i.EventID}},
		event,
	); err != nil {
		return
	}
	database := i.ParticipationDatabase(token, event)
	if err = ParticipationCollection.InsertOne(ctx, database); err != nil {
		return
	}
	//get event by id
	if event, err = EventGetInternalByID(ctx, &models.EventParam{ID: i.EventID}); err != nil {
		return
	}
	eventFilter := bson.D{{Key: "_id", Value: event.ID}}
	var updateEvent bson.D
	if event.TypeOfEvent == "crew_meeting" {
		//if the type of event is crew_meeting, then the value of confirmed and total is increased by one.
		updateEvent = bson.D{{Key: "applications.total", Value: 1}, {Key: "applications.confirmed", Value: 1}}
	} else {
		//else the participation is a request for an asp and so the requested count is increased by one.
		updateEvent = bson.D{{Key: "applications.total", Value: 1}, {Key: "applications.requested", Value: 1}}
	}
	if err = EventCollection.UpdateOne(ctx, eventFilter, vmdb.UpdateInc(updateEvent), nil); err != nil {
		return
	}
	filter := database.Match()
	if err = ParticipationCollection.AggregateOne(
		ctx,
		models.ParticipationPipeline().Match(filter).Pipe,
		&result,
	); err != nil {
		return
	}
	return
}

func ParticipationGet(ctx context.Context, i *models.ParticipationQuery, token *models.AccessToken) (result *[]models.Participation, err error) {
	if err = models.ParticipationPermission(token); err != nil {
		return
	}
	filter := i.PermittedFilter(token)
	result = new([]models.Participation)
	if err = ParticipationCollection.Aggregate(
		ctx,
		models.ParticipationPipeline().Match(filter).Pipe,
		result,
	); err != nil {
		return
	}
	return
}

func ParticipationGetByID(ctx context.Context, i *models.ParticipationParam, token *models.AccessToken) (result *models.Participation, err error) {
	filter := i.PermittedFilter(token)
	if err = ParticipationCollection.AggregateOne(
		ctx,
		models.ParticipationPipeline().Match(filter).Pipe,
		&result,
	); err != nil {
		return
	}
	return
}

func ParticipationUserGet(ctx context.Context, i *models.ParticipationQuery, token *models.AccessToken) (result *[]models.UserParticipation, err error) {
	filter := i.FilterUser(token)
	result = new([]models.UserParticipation)
	if err = ParticipationCollection.Aggregate(
		ctx,
		models.ParticipationPipeline().Match(filter).Pipe,
		result,
	); err != nil {
		return
	}

	return
}

func ParticipationAspGet(ctx context.Context, i *models.ParticipationQuery, token *models.AccessToken) (result *models.EventDetails, err error) {
	filter := i.FilterAspInformation(token)
	participation := new(models.Participation)
	if err = ParticipationCollection.AggregateOne(
		ctx,
		models.ParticipationAspPipeline().Match(filter).Pipe,
		participation,
	); err != nil {
		return
	}
	result = result.FromParticipationEvent(participation.Event)
	return
}

func ParticipationEventGet(ctx context.Context, i *models.EventParam, token *models.AccessToken) (result *[]models.EventParticipation, err error) {
	filter := i.FilterEvent(token)
	result = new([]models.EventParticipation)
	if err = ParticipationCollection.Aggregate(
		ctx,
		models.ParticipationPipeline().Match(filter).Pipe,
		result,
	); err != nil {
		return
	}
	return
}

func ParticipationUpdate(ctx context.Context, i *models.ParticipationUpdate, token *models.AccessToken) (result *models.Participation, err error) {
	participation := new(models.Participation)
	if err = ParticipationCollection.AggregateOne(
		ctx,
		models.ParticipationPipeline().Match(i.Match()).Pipe,
		participation,
	); err != nil {
		return
	}
	if err = i.ParticipationUpdatePermission(token, participation); err != nil {
		return
	}
	filter := i.Match()
	if err = ParticipationCollection.UpdateOneAggregate(
		ctx,
		filter,
		vmdb.UpdateSet(i),
		&result,
		models.ParticipationPipeline().Match(i.Match()).Pipe,
	); err != nil {
		return
	}
	applications := new(models.EventApplications)
	applicationsUpdate := participation.UpdateEventApplicationsUpdate(-1, applications)
	applicationsUpdate = result.UpdateEventApplicationsUpdate(1, &applicationsUpdate.Applications)
	if _, err = EventApplicationsUpdate(ctx, applicationsUpdate); err != nil {
		return
	}
	if participation.Status != result.Status {
		if result.Status == "confirmed" || result.Status == "rejected" {
			ParticipationNotification(ctx, result)
		}
		if result.Status == "withdrawn" {
			ParticipationWithdrawnNotification(ctx, result)
		}
	}
	return
}

func ParticipationDelete(ctx context.Context, i *models.ParticipationParam, token *models.AccessToken) (err error) {
	if err = models.ParticipationDeletePermission(token); err != nil {
		return
	}
	participation := new(models.Participation)
	if participation, err = ParticipationGetByID(ctx, i, token); err != nil {
		return
	}
	applications := new(models.EventApplications)
	applicationsUpdate := participation.UpdateEventApplicationsUpdate(-1, applications)
	if _, err = EventApplicationsUpdate(ctx, applicationsUpdate); err != nil {
		return
	}
	if err = ParticipationCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: i.ID}}); err != nil {
		return
	}
	return
}

func ParticipationNotification(ctx context.Context, i *models.Participation) (err error) {

	template := "participation_confirm"
	if i.Status == "rejected" {
		template = "participation_reject"
	}
	mail := vcago.NewMailData(i.User.Email, "pool-backend", template, "pool", i.User.Country)
	mail.AddUser(i.User.User())
	mail.AddContent(i.ToContent())
	vcago.Nats.Publish("system.mail.job", mail)

	return
}

func ParticipationCreateNotification(ctx context.Context, i *models.Participation) (err error) {
	if i.Event.CrewID == "" || i.Event.EventASPID == "" {
		return vcago.NewNotFound(models.CrewCollection, i)
	}

	template := "participation_create"
	eventAps := new(models.User)
	if eventAps, err = UsersGetByID(ctx, &models.UserParam{ID: i.Event.EventASPID}); err != nil {
		return
	}
	mail := vcago.NewMailData(eventAps.Email, "pool-backend", template, "pool", eventAps.Country)
	mail.AddUser(eventAps.User())
	mail.AddContent(i.ToContent())
	vcago.Nats.Publish("system.mail.job", mail)
	return
}

func ParticipationWithdrawnNotification(ctx context.Context, i *models.Participation) (err error) {

	if i.Event.CrewID == "" || i.Event.EventASPID == "" {
		return vcago.NewNotFound(models.CrewCollection, i)
	}

	template := "participation_withdrawn"
	eventAps := new(models.User)
	if eventAps, err = UsersGetByID(ctx, &models.UserParam{ID: i.Event.EventASPID}); err != nil {
		return
	}
	mail := vcago.NewMailData(eventAps.Email, "pool-backend", template, "pool", eventAps.Country)
	mail.AddUser(eventAps.User())
	mail.AddContent(i.ToContent())
	vcago.Nats.Publish("system.mail.job", mail)
	return
}
