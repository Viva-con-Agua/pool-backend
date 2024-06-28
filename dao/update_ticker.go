package dao

import (
	"context"
	"fmt"
	"log"
	"pool-backend/models"
	"time"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateTicker() {
	ticker := time.NewTicker(1 * time.Hour)
	quit := make(chan struct{})
	EventStateFinishTicker()
	EventStateClosedTicker()
	go func() {
		for {
			select {
			case <-ticker.C:
				EventStateFinishTicker()
				EventStateClosedTicker()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// EventStateFinishTicker sets all published events that are already over to finished
func EventStateFinishTicker() {
	filter := vmdb.NewFilter()
	filter.EqualString("event_state.state", "published")
	filter.LteInt64("end_at", fmt.Sprint(time.Now().Unix()))
	update := bson.D{{Key: "event_state.state", Value: "finished"}}
	events := new([]models.Event)
	if err := EventCollection.Find(context.Background(), filter.Bson(), events); err != nil {
		log.Print(err)
	}
	if err := EventCollection.UpdateMany(context.Background(), filter.Bson(), vmdb.UpdateSet(update)); err != nil {
		log.Print(err)
	}
	for _, value := range *events {
		filterTaking := bson.D{{Key: "_id", Value: value.TakingID}, {Key: "taking_id", Value: bson.D{{Key: "$ne", Value: ""}}}}
		updateTaking := bson.D{{Key: "date_of_taking", Value: value.EndAt}}
		if err := TakingCollection.UpdateOne(context.Background(), filterTaking, vmdb.UpdateSet(updateTaking), nil); err != nil {
			log.Print(err)
		}
	}

}

// EventStateClosed
func EventStateClosedTicker() {
	filter := vmdb.NewFilter()
	confirmedFilter := bson.E{Key: "$or", Value: bson.A{
		bson.D{
			{Key: "state.confirmed.amount", Value: bson.D{{Key: "$gte", Value: 1}}},
			{Key: "state.wait.amount", Value: 0},
			{Key: "state.open.amount", Value: 0},
		},
		bson.D{{Key: "state.no_income", Value: true}},
	}}
	filter.Append(confirmedFilter)
	filter.EqualString("event.event_state.state", "finished")
	pipeline := models.TakingPipeline().Match(filter.Bson()).Pipe
	takings := []models.Taking{}
	if err := TakingCollection.Aggregate(context.Background(), pipeline, takings); err != nil {
		log.Print(err)
	}
	for i := range takings {
		updateFilter := bson.D{{Key: "_id", Value: takings[i].Event.ID}}
		update := bson.D{{Key: "event_state.state", Value: "closed"}}
		e := new(models.Event)
		if err := EventCollection.UpdateOne(context.Background(), updateFilter, vmdb.UpdateSet(update), e); err != nil {
			log.Print(err)
		}
		if err := IDjango.Post(i, "/v1/pool/taking/create/"); err != nil {
			log.Print(err)
		}
		if err := IDjango.Post(e, "/v1/pool/event/update/"); err != nil {
			log.Print(err)
		}
		participations := new([]models.Participation)
		if err := ParticipationCollection.Aggregate(
			context.Background(),
			models.ParticipationPipeline().Match(bson.D{{Key: "event_id", Value: e.ID}}).Pipe,
			participations,
		); err != nil {
			log.Print(err)
		}

		if err := IDjango.Post(participations, "/v1/pool/participations/create/"); err != nil {
			log.Print(err)
		}

	}
}
