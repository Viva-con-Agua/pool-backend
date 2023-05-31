package dao

import (
	"context"
	"pool-backend/models"

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
		models.ParticipationPipeline().Match(filter.Bson()).Pipe,
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
		models.ParticipationAspPipeline().Match(filter.Bson()).Pipe,
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
		models.ParticipationPipeline().Match(filter.Bson()).Pipe,
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
	filter := i.PermittedFilter(token)
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

func ParticipationUpdateStatus(ctx context.Context, i *models.ParticipationStateRequest, token *vcapool.AccessToken) (result *models.Participation, err error) {
	filter := i.PermittedFilter(token)
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
