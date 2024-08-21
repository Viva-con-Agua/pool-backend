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
	UserActiveStateTicker()
	go func() {
		for {
			select {
			case <-ticker.C:
				EventStateFinishTicker()
				EventStateClosedTicker()
				UserActiveStateTicker()
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
	if err := EventCollection.UpdateMany(context.Background(), filter.Bson(), vmdb.UpdateSet(update)); err != nil {
		if !vmdb.ErrNoDocuments(err) {
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
	if err := TakingCollection.Aggregate(context.Background(), pipeline, &takings); err != nil {
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

func UserActiveStateTicker() {
	checkDate := time.Now().Unix() - 15768000
	filter := bson.D{
		{Key: "last_login_date", Value: bson.D{{Key: "$lte", Value: checkDate}}},
		{Key: "active.status", Value: "confirmed"},
	}
	userList := []models.User{}
	pipeline := models.UserPipeline(false).Match(filter).Pipe
	if err := UserCollection.Aggregate(context.Background(), pipeline, &userList); err != nil {
		log.Print(err)
	}
	for _, user := range userList {
		update := bson.D{{Key: "status", Value: "rejected"}}
		userFilter := bson.D{{Key: "_id", Value: user.Active.ID}}
		if err := ActiveCollection.UpdateOne(context.Background(), userFilter, vmdb.UpdateSet(update), nil); err != nil {
			log.Print(err)
		}
	}

}
