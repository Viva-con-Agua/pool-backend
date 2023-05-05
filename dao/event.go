package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func EventInsert(ctx context.Context, i *models.EventDatabase, token *vcapool.AccessToken) (r *models.Event, err error) {
	taking := models.TakingDatabase{
		ID:       uuid.NewString(),
		Name:     i.Name,
		CrewID:   i.CrewID,
		Type:     "automatically",
		Modified: vmod.NewModified(),
	}
	i.TakingID = taking.ID
	if err = EventCollection.InsertOne(ctx, i); err != nil {
		return
	}
	if err = TakingCollection.InsertOne(ctx, taking); err != nil {
		return
	}
	eventActivity := models.NewActivityDB(i.CreatorID, "event", i.ID, "Event created", "created")
	if err = ActivityCollection.InsertOne(ctx, eventActivity); err != nil {
		return
	}
	takingActivity := models.NewActivityDB(i.CreatorID, "taking", taking.ID, "Taking automatically created for Event", "auto_created")
	if err = ActivityCollection.InsertOne(ctx, takingActivity); err != nil {
		return
	}
	r = new(models.Event)
	if err = EventCollection.AggregateOne(ctx, models.EventPipeline(token).Match(i.Match()).Pipe, r); err != nil {
		return
	}
	return
}

func EventGet(ctx context.Context, i *models.EventQuery, token *vcapool.AccessToken) (result *[]models.Event, err error) {
	filter := i.CreateFilter()
	if !token.Roles.Validate("employee;admin") && !token.PoolRoles.Validate("network;operation;education") {
		filter.EqualStringList("event_state.state", []string{"published", "finished", "closed"})
	} else if !token.Roles.Validate("employee;admin") {
		noCrewMatch := vmdb.NewFilter()
		crewMatch := vmdb.NewFilter()
		crewMatch.EqualString("crew_id", token.CrewID)
		noCrewMatch.EqualStringList("event_state.state", []string{"published", "finished", "closed"})
		filter.Append(bson.E{Key: "$or", Value: bson.A{noCrewMatch.Bson(), crewMatch.Bson()}})
	}

	pipeline := models.EventPipeline(token).Match(filter.Bson()).Pipe
	result = new([]models.Event)

	if err = EventCollection.Aggregate(ctx, pipeline, result); err != nil {
		return
	}
	return
}

func EventGetPublic(ctx context.Context, i *models.EventQuery) (result *[]models.Event, err error) {
	i.EventState = []string{"published", "finished", "closed"}
	filter := i.Filter()
	pipeline := models.EventPipelinePublic().Match(filter).Pipe
	result = new([]models.Event)

	if err = EventCollection.Aggregate(ctx, pipeline, result); err != nil {
		return
	}
	return
}

func EventDelete(ctx context.Context, id string) (err error) {
	if err = EventCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}}); err != nil {
		return
	}
	if err = ParticipationCollection.TryDeleteMany(ctx, bson.D{{Key: "event_id", Value: id}}); err != nil {
		return
	}
	return
}


func EventImport(ctx context.Context, import *models.EventImport)