package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func EventInsert(ctx context.Context, i *models.EventDatabase) (r *models.Event, err error) {
	taking := models.TakingDatabase{
		ID:       uuid.NewString(),
		Name:     i.Name,
		CrewID:   i.CrewID,
		Type:     "automatically",
		Status:   "blocked",
		Modified: vmod.NewModified(),
	}
	i.TakingID = taking.ID
	if err = EventCollection.InsertOne(ctx, i); err != nil {
		return
	}
	if err = TakingCollection.InsertOne(ctx, taking); err != nil {
		return
	}
	eventActivity := models.NewActivityDB(i.CreatorID, "event", i.ID, "Event created")
	if err = ActivityCollection.InsertOne(ctx, eventActivity); err != nil {
		return
	}
	takingActivity := models.NewActivityDB(i.CreatorID, "taking", taking.ID, "Taking automatically created for Event")
	if err = ActivityCollection.InsertOne(ctx, takingActivity); err != nil {
		return
	}
	r = new(models.Event)
	if err = EventCollection.AggregateOne(ctx, models.EventPipeline().Match(i.Match()).Pipe, r); err != nil {
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
