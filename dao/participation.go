package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func ParticipationInsert(ctx context.Context, i *models.ParticipationCreate, token *vcapool.AccessToken) (result *models.Participation, err error) {
	database := i.ParticipationDatabase(token)
	if err = ParticipationCollection.InsertOne(ctx, database); err != nil {
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

func ParticipationGet(ctx context.Context, i *models.ParticipationQuery, token *vcapool.AccessToken) (result *[]models.Participation, err error) {
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

func ParticipationGetByID(ctx context.Context, i *models.ParticipationParam, token *vcapool.AccessToken) (result *models.Participation, err error) {
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

func ParticipationUserGet(ctx context.Context, i *models.ParticipationQuery, token *vcapool.AccessToken) (result *[]models.UserParticipation, err error) {
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

func ParticipationAspGet(ctx context.Context, i *models.ParticipationQuery, token *vcapool.AccessToken) (result *models.EventDetails, err error) {
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

func ParticipationEventGet(ctx context.Context, i *models.EventParam, token *vcapool.AccessToken) (result *[]models.EventParticipation, err error) {
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

func ParticipationUpdate(ctx context.Context, i *models.ParticipationUpdate, token *vcapool.AccessToken) (result *models.Participation, err error) {
	event := new(models.Participation)
	if err = ParticipationCollection.AggregateOne(
		ctx,
		models.ParticipationPipeline().Match(i.Match()).Pipe,
		event,
	); err != nil {
		return
	}
	if err = models.ParticipationUpdatePermission(token, event); err != nil {
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
	return
}

func ParticipationDelete(ctx context.Context, i *models.ParticipationParam, token *vcapool.AccessToken) (err error) {
	if err = models.ParticipationDeletePermission(token); err != nil {
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
	mail := vcago.NewMailData(i.User.Email, "pool-backend", template, i.User.Country)
	mail.AddUser(i.User.User())
	mail.AddContent(i.ToContent())
	vcago.Nats.Publish("system.mail.job", mail)
	//notification := vcago.NewMNotificationData(result.User.Email, "pool-backend", template, result.User.Country, token.ID)
	//notification.AddUser(result.User.User())
	//notification.AddContent(result.ToContent())
	//vcago.Nats.Publish("system.notification.job", notification)
	return
}

func ParticipationCreateNotification(ctx context.Context, i *models.Participation) (err error) {

	template := "participation_create"

	users := new([]models.User)
	filter := i.Event.FilterCrew()
	if err = UserCollection.Aggregate(ctx, models.UserPipeline(false).Match(filter).Pipe, users); err != nil {
		return
	}
	eventAps := new(models.User)
	if eventAps, err = UsersGetByID(ctx, &models.UserParam{ID: i.Event.EventASPID}); err != nil {
		return
	}

	mail := vcago.NewMailData(eventAps.Email, "pool-backend", template, eventAps.Country)
	mail.AddUser(eventAps.User())
	mail.AddContent(i.ToContent())
	vcago.Nats.Publish("system.mail.job", mail)
	return
}
