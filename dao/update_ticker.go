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
	EventStateUpdateTicker()
	go func() {
		for {
			select {
			case <-ticker.C:
				EventStateUpdateTicker()
				EventStateNoIncome()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func EventStateUpdateTicker() {
	filter := vmdb.NewFilter()
	filter.EqualString("event_state.state", "published")
	filter.LteInt64("end_at", fmt.Sprint(time.Now().Unix()))
	update := bson.D{{Key: "event_state.state", Value: "finished"}}
	if err := EventCollection.UpdateMany(context.Background(), filter.Bson(), vmdb.UpdateSet(update)); err != nil {
		log.Print(err)
	}
}
func EventStateNoIncome() {
	filter := vmdb.NewFilter()
	filter.EqualBool("no_income", "true")
	filter.EqualString("event.event_state.state", "finished")
	pipeline := models.TakingPipelineTicker().Match(filter.Bson()).Pipe
	takings := []models.Taking{}
	if err := TakingCollection.Aggregate(context.Background(), pipeline, takings); err != nil {
		log.Print(err)
	}
	for i := range takings {
		updateFilter := bson.D{{Key: "_id", Value: takings[i].Event.ID}}
		update := bson.D{{Key: "event_state.state", Value: "closed"}}
		if err := TakingCollection.UpdateOne(context.Background(), updateFilter, vmdb.UpdateSet(update), nil); err != nil {
			log.Print(err)
		}

	}
}
