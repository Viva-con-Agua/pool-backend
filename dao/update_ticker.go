package dao

import (
	"context"
	"fmt"
	"log"
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
